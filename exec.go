package easydb

import (
	"context"
)

func (db *PgxDB) exec() pgxquerier {
	if db.tx != nil {
		return db.tx
	}

	return db.pool
}

func (db *PgxDB) Exec(ctx context.Context, query string, args ...any) (CommandTag, error) {
	query = cleanQuery(query)

	cmd, err := db.exec().Exec(ctx, query, args...)

	return CommandTag{pgxTag: cmd}, err
}

func (db *PgxDB) NamedExec(ctx context.Context, query string, arg any) (CommandTag, error) {
	q, args, err := prepareNamedQuery(query, arg)
	if err != nil {
		return CommandTag{}, err
	}

	cmd, err := db.exec().Exec(ctx, q, args...)

	return CommandTag{pgxTag: cmd}, err
}

func (db *PgxDB) Close() {
	db.pool.Close()
}
