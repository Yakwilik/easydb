package easydb

import (
	"context"
	"fmt"
)

func (db *PgxDB) WithTx(ctx context.Context, fn func(tx Querier) error) (err error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}

	// Оборачиваем в defer логику обработки ошибок и паник
	defer func() {
		if r := recover(); r != nil {
			// Пытаемся откатить транзакцию после паники
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				panic(fmt.Errorf("panic occurred: %v, and rollback failed: %w", r, rollbackErr))
			}
			panic(r)
		}

		if err != nil {
			// Если ошибка во время выполнения функции
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				err = fmt.Errorf("tx rollback failed: %w (original error: %w)", rollbackErr, err)
			}
			return
		}

		// Если всё ок — коммитим
		err = tx.Commit(ctx)
	}()

	err = fn(&PgxDB{pool: db.pool, tx: tx})
	return err
}

func (db *PgxDB) Begin(ctx context.Context) (Tx, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &PgxDB{pool: db.pool, tx: tx}, nil
}

func (db *PgxDB) Commit(ctx context.Context) error {
	return db.tx.Commit(ctx)
}

func (db *PgxDB) Rollback(ctx context.Context) error {
	return db.tx.Rollback(ctx)
}
