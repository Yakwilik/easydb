package easydb

import (
	"context"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
)

type Queryer interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Get(ctx context.Context, dest interface{}, query string, args ...any) error
	Select(ctx context.Context, dest interface{}, query string, args ...any) error

	NamedExec(ctx context.Context, query string, arg interface{}) (pgconn.CommandTag, error)
	NamedGet(ctx context.Context, dest interface{}, query string, arg interface{}) error
	NamedSelect(ctx context.Context, dest interface{}, query string, arg interface{}) error

	GetQuerier() querier
}

type DB interface {
	Queryer

	WithTx(ctx context.Context, fn func(tx Queryer) error) error
	Begin(ctx context.Context) (Tx, error)
}

type Tx interface {
	Queryer
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type pgxDB struct {
	pool *pgxpool.Pool
	tx   pgx.Tx // может быть nil
}

func (db *pgxDB) GetQuerier() querier {
	return db.exec()
}

func New(pool *pgxpool.Pool) *pgxDB {
	return &pgxDB{pool: pool}
}

func (db *pgxDB) exec() querier {
	if db.tx != nil {
		return db.tx
	}
	return db.pool
}

type querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

func (db *pgxDB) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	query = sqlx.Rebind(sqlx.DOLLAR, strings.TrimSpace(query))
	return db.exec().Exec(ctx, query, args...)
}

func (db *pgxDB) Get(ctx context.Context, dest interface{}, query string, args ...any) error {
	query = sqlx.Rebind(sqlx.DOLLAR, strings.TrimSpace(query))

	rows, err := db.exec().Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return pgx.ErrNoRows
	}

	return scanRowToStruct(dest, rows)
}

func (db *pgxDB) Select(ctx context.Context, dest interface{}, query string, args ...any) error {
	query = sqlx.Rebind(sqlx.DOLLAR, strings.TrimSpace(query))

	rows, err := db.exec().Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return scanRowsToSlice(dest, rows)
}

func (db *pgxDB) NamedExec(ctx context.Context, query string, arg interface{}) (pgconn.CommandTag, error) {
	q, args, err := sqlx.Named(query, arg)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	q = sqlx.Rebind(sqlx.DOLLAR, q)
	return db.exec().Exec(ctx, q, args...)
}

func (db *pgxDB) NamedGet(ctx context.Context, dest interface{}, query string, arg interface{}) error {
	q, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}
	q = sqlx.Rebind(sqlx.DOLLAR, q)

	rows, err := db.exec().Query(ctx, q, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return pgx.ErrNoRows
	}

	return scanRowToStruct(dest, rows)
}

func (db *pgxDB) NamedSelect(ctx context.Context, dest interface{}, query string, arg interface{}) error {
	q, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}
	q = sqlx.Rebind(sqlx.DOLLAR, q)

	rows, err := db.exec().Query(ctx, q, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return scanRowsToSlice(dest, rows)
}

func (db *pgxDB) WithTx(ctx context.Context, fn func(tx Queryer) error) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}

	txdb := &pgxDB{pool: db.pool, tx: tx}

	if err := fn(txdb); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func (db *pgxDB) Begin(ctx context.Context) (Tx, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &pgxDB{pool: db.pool, tx: tx}, nil
}

func (db *pgxDB) Commit(ctx context.Context) error {
	return db.tx.Commit(ctx)
}

func (db *pgxDB) Rollback(ctx context.Context) error {
	return db.tx.Rollback(ctx)
}

func scanRowToStruct(dest interface{}, rows pgx.Rows) error {
	return pgxscan.ScanOne(dest, rows)
}

func scanRowsToSlice(dest interface{}, rows pgx.Rows) error {
	return pgxscan.ScanAll(dest, rows)
}
