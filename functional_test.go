//go:build unit

package easydb

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

type user struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

func newMockQuerier(t *testing.T) (pgxmock.PgxPoolIface, *pgxpool.Pool) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	return mock, nil
}

func TestGet(t *testing.T) {
	mock, _ := newMockQuerier(t)
	defer mock.Close()

	rows := pgxmock.NewRows([]string{"id", "name"}).AddRow(1, "Alice")
	mock.ExpectQuery("SELECT id,name FROM users WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	result, err := Get[user](context.Background(), mock, "SELECT id, name FROM users WHERE id = $1", 1)
	require.NoError(t, err)
	require.Equal(t, 1, result.ID)
	require.Equal(t, "Alice", result.Name)
}

func TestSelect(t *testing.T) {
	mock, _ := newMockQuerier(t)
	defer mock.Close()

	rows := pgxmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Alice").
		AddRow(2, "Bob")

	mock.ExpectQuery("SELECT id,name FROM users").
		WillReturnRows(rows)

	result, err := Select[user](context.Background(), mock, "SELECT id, name FROM users")
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, "Bob", result[1].Name)
}

func TestExec(t *testing.T) {
	mock, _ := newMockQuerier(t)
	defer mock.Close()

	mock.ExpectExec("UPDATE users SET name = \\$1 WHERE id = \\$2").
		WithArgs("Alice", 1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	tag, err := Exec(context.Background(), mock, "UPDATE users SET name = $1 WHERE id = $2", "Alice", 1)
	require.NoError(t, err)
	require.True(t, tag.Update())
}

func TestNamedGet(t *testing.T) {
	mock, _ := newMockQuerier(t)
	defer mock.Close()

	rows := pgxmock.NewRows([]string{"id", "name"}).AddRow(1, "Alice")

	mock.ExpectQuery("SELECT id,name FROM users WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	result, err := NamedGet[user](context.Background(), mock, "SELECT id, name FROM users WHERE id = :id", map[string]interface{}{"id": 1})
	require.NoError(t, err)
	require.Equal(t, "Alice", result.Name)
}

func TestNamedSelect(t *testing.T) {
	mock, _ := newMockQuerier(t)
	defer mock.Close()

	rows := pgxmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Alice").
		AddRow(2, "Bob")

	mock.ExpectQuery("SELECT id,name FROM users WHERE name != \\$1").
		WithArgs("Eve").
		WillReturnRows(rows)

	result, err := NamedSelect[user](context.Background(), mock, "SELECT id, name FROM users WHERE name != :name", map[string]interface{}{"name": "Eve"})
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, "Bob", result[1].Name)
}

func TestNamedExec(t *testing.T) {
	mock, _ := newMockQuerier(t)
	defer mock.Close()

	mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
		WithArgs(42).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	tag, err := NamedExec(context.Background(), mock, "DELETE FROM users WHERE id = :id", map[string]interface{}{"id": 42})
	require.NoError(t, err)
	require.True(t, tag.Delete())
}
