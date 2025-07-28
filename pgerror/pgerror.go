package pgerror

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func AsPGError(err error) *pgconn.PgError {
	if pgErr := new(pgconn.PgError); errors.As(err, &pgErr) {
		return pgErr
	}

	return nil
}
