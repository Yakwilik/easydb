package easydb

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

//go:generate go tool mockery --name Querier --filename querier_mock.go --inpackage --with-expecter
type Querier interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Get(ctx context.Context, dest interface{}, query string, args ...any) error
	Select(ctx context.Context, dest interface{}, query string, args ...any) error

	NamedExec(ctx context.Context, query string, arg interface{}) (pgconn.CommandTag, error)
	NamedGet(ctx context.Context, dest interface{}, query string, arg interface{}) error
	NamedSelect(ctx context.Context, dest interface{}, query string, arg interface{}) error

	GetQuerier() pgxquerier
}

//go:generate go tool mockery --name DB --filename db_mock.go --inpackage --with-expecter
type DB interface {
	Querier

	WithTx(ctx context.Context, fn func(tx Querier) error) error
	Begin(ctx context.Context) (Tx, error)

	Close()
}

//go:generate go tool mockery --name Tx --filename querier_mock.go --inpackage --with-expecter
type Tx interface {
	Querier
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

//go:generate go tool mockery --name pgxquerier --filename pgxquerier_mock.go --inpackage --with-expecter
type pgxquerier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}
