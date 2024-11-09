package main

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxbolgarin/errm"
)

var (
	ErrEmptyResult     = errm.New("result is empty")
	ErrUniqueViolation = errm.New("unique violation")
	ErrNotAffected     = errm.New("there is no affected rows")
)

type DBConfig struct {
	Address  string `env:"DB_ADDRESS"`
	DBName   string `env:"DB_NAME"`
	Username string `env:"DB_USER"`
	Password string `env:"DB_PASS"`
}

func (cfg DBConfig) BuildURL() string {
	if cfg.Address == "" {
		cfg.Address = "127.0.0.1:5432"
	}

	var url strings.Builder

	url.WriteString("postgres://")
	if cfg.Username != "" && cfg.Password != "" {
		url.WriteString(cfg.Username + ":" + cfg.Password + "@")
	}

	url.WriteString(cfg.Address)
	if cfg.DBName != "" {
		url.WriteString("/" + cfg.DBName)
	}

	return url.String()
}

var DefaultScanAPI = mustNewAPI(dbscan.WithAllowUnknownColumns(true), dbscan.WithStructTagKey("db"))

type DB struct {
	*pgxpool.Pool
}

func NewDB(ctx context.Context) (*DB, error) {
	var cfg DBConfig
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, errm.Wrap(err, "read environment variables")
	}

	conn, err := pgxpool.New(ctx, cfg.BuildURL())
	if err != nil {
		return nil, errm.Wrap(err, "*DBect to database")
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, errm.Wrap(err, "ping database")
	}

	db := &DB{
		conn,
	}

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	log.Println("connected to database at "+cfg.Address, cfg.DBName)

	return db, nil
}

func Exec(ctx context.Context, c *DB, query string, args ...any) (int64, error) {
	tag, err := c.Exec(ctx, query, args...)
	if err != nil {
		return 0, handleErr(err, "exec")
	}
	return handleTag(tag)
}

func QueryOne[T any](ctx context.Context, c *DB, query string, args ...any) (out T, err error) {
	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		return out, handleErr(err, "query")
	}
	if err := DefaultScanAPI.ScanOne(&out, rows); err != nil {
		return out, handleErr(err, "scan one")
	}
	return out, nil
}

func QueryAll[T any](ctx context.Context, c *DB, query string, args ...any) (out []T, err error) {
	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		return nil, handleErr(err, "query")
	}
	if err := DefaultScanAPI.ScanAll(&out, rows); err != nil {
		return nil, handleErr(err, "scan all")
	}
	return out, nil
}

func handleTag(tag pgconn.CommandTag) (int64, error) {
	rows := tag.RowsAffected()
	if rows == 0 {
		switch {
		case tag.Insert():
			return 0, ErrUniqueViolation
		case tag.Update(), tag.Delete():
			return 0, ErrNotAffected
		}
	}
	return rows, nil
}

func handleErr(err error, msg string) error {
	var e *pgconn.PgError

	// https://www.postgresql.org/docs/16/errcodes-appendix.html
	if errors.As(err, &e) {
		if e.Code == pgerrcode.UniqueViolation {
			return errm.Wrap(ErrUniqueViolation, msg)
		}
		return errm.Wrap(err, msg, "code", e.Code)
	}

	if errm.Is(err, pgx.ErrNoRows) {
		return errm.Wrap(ErrEmptyResult, msg)
	}

	if dbscan.NotFound(err) {
		return errm.Wrap(ErrEmptyResult, msg)
	}

	return errm.Wrap(err, msg)
}

func mustNewAPI(opts ...dbscan.APIOption) *pgxscan.API {
	dbScanAPI, err := pgxscan.NewDBScanAPI(opts...)
	if err != nil {
		panic("new dbscan api: " + err.Error())
	}
	api, err := pgxscan.NewAPI(dbScanAPI)
	if err != nil {
		panic("new pgxscan api: " + err.Error())
	}
	return api
}
