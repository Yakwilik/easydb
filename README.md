# easydb

Лёгкая обёртка над pgx (PostgreSQL) с удобными generics-хелперами, поддержкой именованных параметров и простым интерфейсом для транзакций.

- **Поддержка generics-хелперов**: `Get[T]`, `Select[T]`, `NamedGet[T]`, `NamedSelect[T]`, `Exec`, `NamedExec`.
- **Интерфейсы для абстракции БД**: `Querier`, `DB`, `Tx`.
- **Управление транзакциями**: `WithTx`, `Begin`, `Commit`, `Rollback`.
- **Конфиг для подключения**: `PGConfig`, фабрики `New`, `NewWithDSN`, `NewPgxDB`.
- **Утилита ошибок**: `pgerror.AsPGError` для распаковки `*pgconn.PgError`.
- **CommandTag**-обёртка над `pgconn.CommandTag`.

---

## Установка

```bash
go get github.com/khasbulatabdullin/easydb
```

Требования: Go 1.20+, PostgreSQL, pgx/v5.

---

## Все способы инициализации

### 1) Из существующего пула `pgxpool.Pool`
```go
pool, err := pgxpool.New(ctx, dsn)
if err != nil { /* ... */ }
db := easydb.New(pool)
defer db.Close()
```

### 2) По готовому DSN (создание пула + Ping)
```go
db, err := easydb.NewWithDSN(ctx, "postgres://user:pass@localhost:5432/app?sslmode=disable")
if err != nil { /* ... */ }
defer db.Close()
```

### 3) Через `PGConfig` (сборка DSN + создание пула + Ping)
```go
cfg := easydb.NewPGConfig().
	WithHost("localhost").
	WithPort("5432").
	WithUser("user").
	WithPassword("pass").
	WithDatabase("app").
	WithSSLMode("disable").
	WithMaxConns(10).
	WithAppName("easydb-svc").
	WithSchema("public")

db, err := easydb.NewPgxDB(ctx, cfg)
if err != nil { /* ... */ }
defer db.Close()
```

### 4) С пользовательскими параметрами в `Params`
```go
cfg := easydb.NewPGConfig()
cfg.Params["pool_min_conns"] = "2"
cfg.Params["connect_timeout"] = "5"
db, err := easydb.NewPgxDB(ctx, cfg)
```

### 5) С контекстом-таймаутом на инициализацию
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

db, err := easydb.NewWithDSN(ctx, dsn)
```

---

## Быстрый старт

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	easydb "github.com/khasbulatabdullin/easydb"
)

type User struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func main() {
	ctx := context.Background()

	pool, _ := pgxpool.New(ctx, "postgres://user:pass@localhost:5432/app?sslmode=disable")
	db := easydb.New(pool)
	defer db.Close()

	// таймаут на запрос
	reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// простые запросы
	var u User
	_ = db.Get(reqCtx, &u, "SELECT id, name FROM users WHERE id=$1", 1)

	var users []User
	_ = db.Select(reqCtx, &users, "SELECT id, name FROM users")

	// именованные параметры
	_, _ = db.NamedExec(reqCtx, "INSERT INTO users(name) VALUES(:name)", map[string]any{"name": "Alice"})

	// транзакция
	_ = db.WithTx(reqCtx, func(tx easydb.Querier) error {
		_, err := tx.Exec(reqCtx, "UPDATE users SET name=$1 WHERE id=$2", "Bob", 1)
		return err
	})

	fmt.Println("ok")
}
```

---

## Интерфейсы и типы

### `type Querier interface`
- `Ping(ctx context.Context) error`
- `Exec(ctx context.Context, query string, args ...any) (CommandTag, error)`
- `Get(ctx context.Context, dest any, query string, args ...any) error`
- `Select(ctx context.Context, dest any, query string, args ...any) error`
- `NamedExec(ctx context.Context, query string, arg any) (CommandTag, error)`
- `NamedGet(ctx context.Context, dest any, query string, arg any) error`
- `NamedSelect(ctx context.Context, dest any, query string, arg any) error`
- `GetQuerier() pgxquerier`

### `type DB interface`
- `WithTx(ctx context.Context, fn func(tx Querier) error) error`
- `Begin(ctx context.Context) (Tx, error)`
- `Close()`

### `type Tx interface`
- `Commit(ctx context.Context) error`
- `Rollback(ctx context.Context) error`

### `type CommandTag`
- `NewCommandTag(tag pgconn.CommandTag) CommandTag`
- `RowsAffected() int64`
- `Insert() bool`, `Update() bool`, `Delete() bool`, `String() string`

### `type PGConfig`
- Поля: `Host`, `Port`, `Username`, `Password`, `DBName`, `Params map[string]string`
- `NewPGConfig()` и флюенты: `WithHost`, `WithPort`, `WithUser`, `WithPassword`, `WithDatabase`, `WithSSLMode`, `WithMaxConns`, `WithAppName`, `WithSchema`

### Конструкторы БД
- `New(pool *pgxpool.Pool) *PgxDB`
- `NewWithDSN(ctx context.Context, dsn string) (*PgxDB, error)`
- `NewPgxDB(ctx context.Context, cfg PGConfig) (*PgxDB, error)`

---

## Расширенные примеры

### Exec и CommandTag
```go
cmd, err := db.Exec(ctx, "UPDATE users SET name=$1 WHERE id=$2", "Alice", 1)
if err != nil { /* ... */ }
fmt.Println(cmd.RowsAffected())
if cmd.Update() { /* обновление */ }
```

### Get/Select в структуры и примитивы
```go
// структура
var u User
_ = db.Get(ctx, &u, "SELECT id, name FROM users WHERE id=$1", 1)

// срез структур
var users []User
_ = db.Select(ctx, &users, "SELECT id, name FROM users WHERE id > $1", 0)

// примитив
var cnt int
_ = db.Get(ctx, &cnt, "SELECT COUNT(*) FROM users")
```

### Именованные параметры: map, struct, раскрытие слайсов
```go
// map[string]any
_, _ = db.NamedExec(ctx, "INSERT INTO users(name, age) VALUES(:name, :age)", map[string]any{"name": "Bob", "age": 42})

// struct
type CreateUser struct { Name string `db:"name"`; Age int `db:"age"` }
_, _ = db.NamedExec(ctx, "INSERT INTO users(name, age) VALUES(:name, :age)", CreateUser{Name: "Eve", Age: 27})

// slice в IN (через sqlx.In внутри)
rows, _ := easydb.NamedSelect[User](ctx, db.GetQuerier(),
	"SELECT id, name FROM users WHERE id IN (:ids)",
	map[string]any{"ids": []int64{1,2,3}},
)
```

### Использование generics-хелперов без `*PgxDB`
```go
// работает с любым исполнителем, реализующим Query/Exec из pgx (pool или tx)
u, err := easydb.Get[User](ctx, db.GetQuerier(), "SELECT id, name FROM users WHERE id=$1", 1)
list, err := easydb.Select[User](ctx, db.GetQuerier(), "SELECT id, name FROM users")
_, err = easydb.Exec(ctx, db.GetQuerier(), "DELETE FROM users WHERE id=$1", 10)
```

### Транзакции: шаблон с `WithTx`
```go
err := db.WithTx(ctx, func(tx easydb.Querier) error {
	if _, err := tx.Exec(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id=$2", 100, 1); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id=$2", 100, 2); err != nil {
		return err
	}
	return nil
})
```

### Транзакции: ручное Begin/Commit/Rollback
```go
trx, err := db.Begin(ctx)
if err != nil { /* ... */ }

var opErr error

defer func() {
	if opErr != nil {
		_ = trx.Rollback(ctx)
		return
	}
	_ = trx.Commit(ctx)
}()

_, opErr = trx.Exec(ctx, "INSERT INTO audit(event) VALUES($1)", "user_created")
if opErr != nil { return }

var total int
opErr = trx.Get(ctx, &total, "SELECT COUNT(*) FROM audit")
```

### Обработка unique_violation в транзакции
```go
import pgxerr "github.com/khasbulatabdullin/easydb/pgerror"

_ = db.WithTx(ctx, func(tx easydb.Querier) error {
	_, err := tx.Exec(ctx, "INSERT INTO users(email) VALUES($1)", "dup@example.com")
	if err != nil {
		if pg := pgxerr.AsPGError(err); pg != nil && pg.Code == "23505" { // unique_violation
			// игнорируем или преобразуем в доменную ошибку
			return nil
		}
		return err
	}
	return nil
})
```

### Контекст и таймауты
```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
var users []User
err := db.Select(ctx, &users, "SELECT id, name FROM users")
```

### Обработка ошибок PostgreSQL
```go
import pgxerr "github.com/khasbulatabdullin/easydb/pgerror"

if err != nil {
	if pg := pgxerr.AsPGError(err); pg != nil {
		if pg.Code == "23505" { /* unique_violation */ }
	}
}
```

### Тестирование с моками (`MockQuerier`)
```go
m := easydb.NewMockQuerier(t)
m.EXPECT().Exec(mock.Anything, "UPDATE users SET name=$1 WHERE id=$2", "Alice", 1).
	Return(easydb.NewCommandTag(pgconn.CommandTag("UPDATE 1")), nil)

cmd, err := m.Exec(context.Background(), "UPDATE users SET name=$1 WHERE id=$2", "Alice", 1)
```

---

## Примеры из тестов

Ниже — выдержки, повторяющие реальные сценарии из `integration_test.go` и `functional_test.go`.

### Integration: выборка и получение записи
```go
// SELECT * ORDER BY id
var users []User
err := db.Select(context.Background(), &users, `SELECT * FROM users ORDER BY id`)
// проверка len(users) >= 2

// SELECT * WHERE id=$1
var user User
err = db.Get(context.Background(), &user, `SELECT * FROM users WHERE id = $1`, 1)
```

### Integration: именованные запросы (map и слайсы)
```go
// NamedSelect по name
var users []User
err := db.NamedSelect(context.Background(), &users, `SELECT * FROM users WHERE name = :name`, map[string]any{
	"name": "Alice",
})

// NamedGet по name
var user User
err = db.NamedGet(context.Background(), &user, `SELECT * FROM users WHERE name = :name`, map[string]any{
	"name": "Bob",
})

// NamedExec вставка одного пользователя
_, err = db.NamedExec(context.Background(), `
	INSERT INTO users (name, email) VALUES (:name, :email)
`, map[string]any{
	"email": "charlie@example.com",
	"name":  "Charlie",
})

// NamedExec вставка многих пользователей из среза структур
usersToInsert := []User{
	{Name: "Dave", Email: "dave@example.com"},
	{Name: "Eve",  Email: "eve@example.com"},
}
_, err = db.NamedExec(context.Background(), `
	INSERT INTO users (name, email) VALUES (:name, :email)
`, usersToInsert)

// NamedSelect с IN (:names) — раскрытие слайса
var inserted []User
err = db.NamedSelect(context.Background(), &inserted, `
	SELECT * FROM users WHERE name IN (:names)
`, map[string]any{
	"names": []string{"Dave", "Eve"},
})
```

### Integration: транзакции `WithTx`
```go
var countBefore, countAfter int
err := db.WithTx(context.Background(), func(tx easydb.Querier) error {
	if err := tx.Get(context.Background(), &countBefore, `SELECT COUNT(*) FROM users`); err != nil {
		return err
	}
	if _, err := tx.Exec(context.Background(), `INSERT INTO users (name, email) VALUES ($1, $2)`, "TestUser", "test@example.com"); err != nil {
		return err
	}
	return tx.Get(context.Background(), &countAfter, `SELECT COUNT(*) FROM users`)
})
```

### Unit: generics-хелперы с pgxmock
```go
// Get[T]
result, err := easydb.Get[user](context.Background(), mock, "SELECT id, name FROM users WHERE id = $1", 1)

// Select[T]
list, err := easydb.Select[user](context.Background(), mock, "SELECT id, name FROM users")

// Exec
cmd, err := easydb.Exec(context.Background(), mock, "UPDATE users SET name = $1 WHERE id = $2", "Alice", 1)
_ = cmd.Update()

// NamedGet
u, err := easydb.NamedGet[user](context.Background(), mock, "SELECT id, name FROM users WHERE id = :id", map[string]any{"id": 1})

// NamedSelect
rows, err := easydb.NamedSelect[user](context.Background(), mock, "SELECT id, name FROM users WHERE name != :name", map[string]any{"name": "Eve"})

// NamedExec
cmd, err = easydb.NamedExec(context.Background(), mock, "DELETE FROM users WHERE id = :id", map[string]any{"id": 42})
_ = cmd.Delete()
```

---

## Короткие сниппеты (для индексаторов)

Ниже — изолированные, короткие примеры. Каждый сниппет самодостаточен.

### Get одной записи
```go
var u User
err := db.Get(ctx, &u, "SELECT id, name FROM users WHERE id=$1", 1)
```

### Select списка
```go
var users []User
err := db.Select(ctx, &users, "SELECT id, name FROM users ORDER BY id")
```

### Exec и проверка CommandTag
```go
cmd, err := db.Exec(ctx, "UPDATE users SET name=$1 WHERE id=$2", "Alice", 1)
_ = cmd.Update()
_ = cmd.RowsAffected()
```

### NamedExec с map
```go
_, err := db.NamedExec(ctx, "INSERT INTO users(name, email) VALUES(:name, :email)", map[string]any{"name": "Bob", "email": "b@example.com"})
```

### NamedGet с map
```go
var u User
err := db.NamedGet(ctx, &u, "SELECT id, name FROM users WHERE name=:name", map[string]any{"name": "Alice"})
```

### NamedSelect со слайсом (IN)
```go
var users []User
err := db.NamedSelect(ctx, &users, "SELECT id, name FROM users WHERE id IN (:ids)", map[string]any{"ids": []int64{1,2,3}})
```

### Generics: Get[T]
```go
u, err := easydb.Get[User](ctx, db.GetQuerier(), "SELECT id, name FROM users WHERE id=$1", 1)
```

### Generics: Select[T]
```go
rows, err := easydb.Select[User](ctx, db.GetQuerier(), "SELECT id, name FROM users")
```

### Generics: Exec
```go
cmd, err := easydb.Exec(ctx, db.GetQuerier(), "DELETE FROM users WHERE id=$1", 42)
_ = cmd.Delete()
```

### Generics: NamedGet[T]
```go
u, err := easydb.NamedGet[User](ctx, db.GetQuerier(), "SELECT id, name FROM users WHERE id = :id", map[string]any{"id": 1})
```

### Generics: NamedSelect[T]
```go
list, err := easydb.NamedSelect[User](ctx, db.GetQuerier(), "SELECT id, name FROM users WHERE name != :name", map[string]any{"name": "Eve"})
```

### Generics: NamedExec
```go
cmd, err := easydb.NamedExec(ctx, db.GetQuerier(), "UPDATE users SET name=:name WHERE id=:id", map[string]any{"id": 1, "name": "Alice"})
_ = cmd.Update()
```

### WithTx обёртка
```go
err := db.WithTx(ctx, func(tx easydb.Querier) error {
	_, err := tx.Exec(ctx, "INSERT INTO audit(event) VALUES($1)", "x")
	return err
})
```

### Ручная транзакция
```go
trx, err := db.Begin(ctx)
if err != nil { /* ... */ }

defer func() {
	if err != nil { _ = trx.Rollback(ctx); return }
	_ = trx.Commit(ctx)
}()

_, err = trx.Exec(ctx, "INSERT INTO audit(event) VALUES($1)", "x")
```

### Обработка pgerror.AsPGError (unique_violation)
```go
import pgxerr "github.com/khasbulatabdullin/easydb/pgerror"

_, err := db.Exec(ctx, "INSERT INTO users(email) VALUES($1)", "dup@example.com")
if pg := pgxerr.AsPGError(err); pg != nil && pg.Code == "23505" {
	// unique_violation
}
```

### Инициализация: NewWithDSN
```go
db, err := easydb.NewWithDSN(ctx, "postgres://user:pass@localhost:5432/app?sslmode=disable")
```

### Инициализация: NewPgxDB + PGConfig
```go
cfg := easydb.NewPGConfig().WithHost("localhost").WithPort("5432").WithUser("user").WithPassword("pass").WithDatabase("app").WithSSLMode("disable")
db, err := easydb.NewPgxDB(ctx, cfg)
```

### Инициализация: New(pgxpool.Pool)
```go
pool, _ := pgxpool.New(ctx, dsn)
db := easydb.New(pool)
```

---

## Лицензия

См. файл `LICENSE`.