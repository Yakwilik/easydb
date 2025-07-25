package easydb

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

// Get выполняет запрос и возвращает один объект T
func Get[T any](ctx context.Context, q querier, sql string, args ...any) (T, error) {
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		var zero T
		return zero, err
	}
	defer rows.Close()

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
}

// Select выполняет запрос и возвращает []T
func Select[T any](ctx context.Context, q querier, sql string, args ...any) ([]T, error) {
	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}

// Exec просто выполняет запрос
func Exec(ctx context.Context, q querier, sql string, args ...any) (pgconn.CommandTag, error) {
	return q.Exec(ctx, sql, args...)
}

// NamedGet выполняет именованный запрос и возвращает один объект T
func NamedGet[T any](ctx context.Context, q querier, query string, arg interface{}) (T, error) {
	named, args, err := sqlx.Named(query, arg)
	if err != nil {
		var zero T
		return zero, err
	}
	named = sqlx.Rebind(sqlx.DOLLAR, named)

	rows, err := q.Query(ctx, named, args...)
	if err != nil {
		var zero T
		return zero, err
	}
	defer rows.Close()

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
}

// NamedSelect выполняет именованный запрос и возвращает []T
func NamedSelect[T any](ctx context.Context, q querier, query string, arg interface{}) ([]T, error) {
	named, args, err := sqlx.Named(query, arg)
	if err != nil {
		return nil, err
	}
	named = sqlx.Rebind(sqlx.DOLLAR, named)

	rows, err := q.Query(ctx, named, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}

// NamedExec выполняет именованный запрос без возврата результата
func NamedExec(ctx context.Context, q querier, query string, arg interface{}) (pgconn.CommandTag, error) {
	named, args, err := sqlx.Named(query, arg)
	if err != nil {
		return pgconn.CommandTag{}, err
	}
	named = sqlx.Rebind(sqlx.DOLLAR, named)

	return q.Exec(ctx, named, args...)
}
