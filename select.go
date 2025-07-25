package easydb

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func (db *PgxDB) Select(ctx context.Context, dest interface{}, query string, args ...any) error {
	query = cleanQuery(query)

	rows, err := db.exec().Query(ctx, query, args...)
	if err != nil {
		return err
	}

	return scanRowsToSlice(dest, rows)
}

func (db *PgxDB) NamedSelect(ctx context.Context, dest interface{}, query string, arg interface{}) error {
	q, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}
	q = cleanQuery(q)

	rows, err := db.exec().Query(ctx, q, args...)
	if err != nil {
		return err
	}

	return scanRowsToSlice(dest, rows)
}
