package store

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")
var ErrInvalidInput = errors.New("invalid input")
var ErrForeignKey = errors.New("foreign key violation")

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
	case "23503": // foreign_key_violation
		return ErrForeignKey
	case "23514": // check_violation
		return ErrInvalidInput
	case "22P02", "22007": // invalid_text_representation, invalid_datetime_format
		return ErrInvalidInput
	case "55P03": // lock_not_available
		return fmt.Errorf("resource locked: %w", ErrConflict)
	case "40001": // serialization_failure
		return fmt.Errorf("serialization failure: %w", ErrConflict)
	case "40P01": // deadlock_detected
		return fmt.Errorf("deadlock detected: %w", ErrConflict)
	}
	return err
}

func Open() (*sql.DB, error) {
	dbHost := os.Getenv("POSTGRES_URL")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	//TODO: don't read password from env var
	dbPassword := os.Getenv("POSTGRES_PASSWORD")

	db, err := sql.Open("pgx", fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort))
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	return db, nil
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

// TODO: move this to pipeline or K8s job
func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}
