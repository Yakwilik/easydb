package easydb

import (
	"context"
)

func (db *PgxDB) Get(ctx context.Context, dest any, query string, args ...any) error {
	query = cleanQuery(query)
	rows, err := db.exec().Query(ctx, query, args...)
	if err != nil {
		return err
	}

	return scanRowToStruct(dest, rows)
}

func (db *PgxDB) NamedGet(ctx context.Context, dest any, query string, arg any) error {
	q, args, err := prepareNamedQuery(query, arg)
	if err != nil {
		return err
	}

	rows, err := db.exec().Query(ctx, q, args...)
	if err != nil {
		return err
	}

	return scanRowToStruct(dest, rows)
}
