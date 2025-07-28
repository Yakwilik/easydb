//go:build integration

package easydb

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

type User struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func setupTestDB(t *testing.T) DB {
	t.Helper()

	dsn := "postgres://testuser:testpass@localhost:5433/testdb"

	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)

	err = pool.Ping(context.Background())
	require.NoError(t, err)

	return New(pool)
}

func TestSelectAndGet(t *testing.T) {
	db := setupTestDB(t)

	t.Run("Select users", func(t *testing.T) {
		var users []User
		err := db.Select(context.Background(), &users, `SELECT * FROM users ORDER BY id`)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(users), 2)
	})

	t.Run("Get user by id", func(t *testing.T) {
		var user User
		err := db.Get(context.Background(), &user, `SELECT * FROM users WHERE id = $1`, 1)
		require.NoError(t, err)
		require.Equal(t, int64(1), user.ID)
	})
}

func TestNamedQueries(t *testing.T) {
	db := setupTestDB(t)

	t.Run("NamedSelect users", func(t *testing.T) {
		var users []User
		err := db.NamedSelect(context.Background(), &users, `SELECT * FROM users WHERE name = :name`, map[string]any{
			"name": "Alice",
		})
		require.NoError(t, err)
		require.NotEmpty(t, users)
	})

	t.Run("NamedGet user", func(t *testing.T) {
		var user User
		err := db.NamedGet(context.Background(), &user, `SELECT * FROM users WHERE name = :name`, map[string]any{
			"name": "Bob",
		})
		require.NoError(t, err)
		require.Equal(t, "Bob", user.Name)
	})

	t.Run("NamedExec insert user", func(t *testing.T) {
		_, err := db.NamedExec(context.Background(), `
			INSERT INTO users (name, email) VALUES (:name, :email)
		`, map[string]any{
			"email": "charlie@example.com",
			"name":  "Charlie",
		})
		require.NoError(t, err)

		// Проверка что пользователь реально вставлен
		var user User
		err = db.NamedGet(context.Background(), &user, `SELECT * FROM users WHERE name = :name`, map[string]any{
			"name": "Charlie",
		})
		require.NoError(t, err)
		require.Equal(t, "Charlie", user.Name)
	})
}

func TestWithTx(t *testing.T) {
	db := setupTestDB(t)

	var countBefore, countAfter int

	err := db.WithTx(context.Background(), func(tx Querier) error {
		var err error

		err = tx.Get(context.Background(), &countBefore, `SELECT COUNT(*) FROM users`)
		if err != nil {
			return err
		}

		_, err = tx.Exec(context.Background(), `INSERT INTO users (name, email) VALUES ($1, $2)`, "TestUser", "test@example.com")
		if err != nil {
			return err
		}

		err = tx.Get(context.Background(), &countAfter, `SELECT COUNT(*) FROM users`)
		return err
	})

	require.NoError(t, err)
	require.Equal(t, countBefore+1, countAfter)
}
