package easydb

import (
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

func scanRowToStruct(dest any, rows pgx.Rows) error {
	return pgxscan.ScanOne(dest, rows)
}

func scanRowsToSlice(dest any, rows pgx.Rows) error {
	return pgxscan.ScanAll(dest, rows)
}
