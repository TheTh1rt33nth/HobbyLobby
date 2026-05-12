package store

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"

	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")
var ErrInvalidInput = errors.New("invalid input")

// translatePgError maps known postgres error codes to store sentinel errors.
// Returns the original error unchanged if it is not a recognised postgres error.
func translatePgError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	switch pgErr.Code {
	case "23505": // unique_violation
		return ErrConflict
	case "22P02", "22007": // invalid_text_representation, invalid_datetime_format
		return ErrInvalidInput
	}
	return err
}

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")

	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	return db, err
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("Migrate: failed to set goose dialect: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("Migrate: failed to run goose up: %w", err)
	}

	return nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}
