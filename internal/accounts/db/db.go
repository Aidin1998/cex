package db

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// OpenAndMigrate opens *sql.DB and runs Goose migrations under db/accounts/migration
func OpenAndMigrate(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	goose.SetDialect("postgres")
	if err := goose.Up(db, "internal/accounts/db/migration"); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
