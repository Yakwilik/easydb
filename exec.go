package easydb

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

func (db *PgxDB) exec() pgxquerier {
	if db.tx != nil {
		return db.tx
	}
	return db.pool
}

func (db *PgxDB) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	query = cleanQuery(query)
	return db.exec().Exec(ctx, query, args...)
}

func (db *PgxDB) NamedExec(ctx context.Context, query string, arg interface{}) (pgconn.CommandTag, error) {
	q, args, err := sqlx.Named(query, arg)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	q = cleanQuery(query)
	return db.exec().Exec(ctx, q, args...)
}

func (db *PgxDB) Close() {
	db.pool.Close()
}
