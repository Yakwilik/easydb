package easydb

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
)

var _ DB = &PgxDB{}
var _ Tx = &PgxDB{}

type PgxDB struct {
	pool *pgxpool.Pool
	tx   pgx.Tx // может быть nil
}

func (db *PgxDB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

func (db *PgxDB) GetQuerier() pgxquerier {
	return db.exec()
}

func New(pool *pgxpool.Pool) *PgxDB {
	return &PgxDB{pool: pool}
}

func NewWithDSN(ctx context.Context, dsn string) (*PgxDB, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PgxDB{pool: pool}, nil
}

func NewPgxDB(ctx context.Context, cfg PGConfig) (*PgxDB, error) {
	// Собираем query-параметры
	values := url.Values{}
	for k, v := range cfg.Params {
		values.Set(k, v)
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		url.QueryEscape(cfg.Username),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		values.Encode(),
	)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &PgxDB{pool: pool}, nil
}

func cleanQuery(query string) string {
	return sqlx.Rebind(sqlx.DOLLAR, replaceSQLChars(query))
}

var rxSpaces = regexp.MustCompile(`\s+`)

// replaceSQLChars remove all non-printable symbols (as a \t, \n) for more comfortable reading SQL ("humanity").
// Need for correct logging database library output (now library cut off multiline query that`s why
// query present in log message only part not all)
func replaceSQLChars(sql string) string {
	sql = strings.Replace(sql, "\n", " ", -1)
	sql = strings.Replace(sql, "\t", " ", -1)
	sql = strings.Replace(sql, " ,", ",", -1)
	sql = strings.Replace(sql, ", ", ",", -1)

	return strings.TrimSpace(rxSpaces.ReplaceAllString(sql, " "))
}

// prepareNamedQuery подготавливает запрос с именованными параметрами:
// - sqlx.Named
// - sqlx.In (если нужно раскрытие слайсов)
// - sqlx.Rebind под `$1, $2, ...`
func prepareNamedQuery(query string, arg any) (string, []any, error) {
	namedQuery, args, err := sqlx.Named(query, arg)
	if err != nil {
		return "", nil, err
	}

	// безопасно даже если нет слайсов
	inQuery, inArgs, err := sqlx.In(namedQuery, args...)
	if err != nil {
		return "", nil, err
	}

	finalQuery := cleanQuery(inQuery)

	return finalQuery, inArgs, nil
}
