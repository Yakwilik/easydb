package easydb

import (
	"github.com/jackc/pgx/v5/pgconn"
)

// CommandTag — обёртка над pgconn.CommandTag, с возможностью расширения
type CommandTag struct {
	pgxTag pgconn.CommandTag
}

// NewCommandTag оборачивает pgconn.CommandTag
func NewCommandTag(tag pgconn.CommandTag) CommandTag {
	return CommandTag{pgxTag: tag}
}

// RowsAffected возвращает число затронутых строк
func (t CommandTag) RowsAffected() int64 {
	return t.pgxTag.RowsAffected()
}

// Insert/Update/Delete и прочие хелперы
func (t CommandTag) Insert() bool   { return t.pgxTag.Insert() }
func (t CommandTag) Update() bool   { return t.pgxTag.Update() }
func (t CommandTag) Delete() bool   { return t.pgxTag.Delete() }
func (t CommandTag) String() string { return t.pgxTag.String() }
