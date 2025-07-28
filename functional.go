package easydb

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Get выполняет запрос и возвращает один объект T
func Get[T any](ctx context.Context, q pgxquerier, sql string, args ...any) (T, error) {
	rows, err := q.Query(ctx, cleanQuery(sql), args...)
	if err != nil {
		var zero T
		return zero, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
}

// Select выполняет запрос и возвращает []T
func Select[T any](ctx context.Context, q pgxquerier, sql string, args ...any) ([]T, error) {
	rows, err := q.Query(ctx, cleanQuery(sql), args...)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}

// Exec просто выполняет запрос
func Exec(ctx context.Context, q pgxquerier, sql string, args ...any) (CommandTag, error) {
	cmdTag, err := q.Exec(ctx, cleanQuery(sql), args...)

	return CommandTag{pgxTag: cmdTag}, err
}

// NamedGet выполняет именованный запрос и возвращает один объект T
func NamedGet[T any](ctx context.Context, q pgxquerier, query string, arg any) (T, error) {
	named, args, err := prepareNamedQuery(query, arg)
	if err != nil {
		var zero T
		return zero, err
	}

	rows, err := q.Query(ctx, named, args...)
	if err != nil {
		var zero T
		return zero, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[T])
}

// NamedSelect выполняет именованный запрос и возвращает []T
func NamedSelect[T any](ctx context.Context, q pgxquerier, query string, arg any) ([]T, error) {
	named, args, err := prepareNamedQuery(query, arg)
	if err != nil {
		return nil, err
	}

	rows, err := q.Query(ctx, named, args...)
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}

// NamedExec выполняет именованный запрос без возврата результата
func NamedExec(ctx context.Context, q pgxquerier, query string, arg any) (CommandTag, error) {
	named, args, err := prepareNamedQuery(query, arg)
	if err != nil {
		return CommandTag{}, err
	}

	cmdTag, err := q.Exec(ctx, named, args...)

	return CommandTag{pgxTag: cmdTag}, err
}
