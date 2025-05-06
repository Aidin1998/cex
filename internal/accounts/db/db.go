package db

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.opentelemetry.io/otel"
)

// NewDB opens a *sql.DB using the provided DSN.
func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// optional: db.SetMaxOpenConns(...)
	return db, nil
}

// ConnectAndMigrate opens the DB, pings it, runs all goose migrations, and returns the live *sql.DB.
func ConnectAndMigrate(ctx context.Context, dsn string) (*sql.DB, error) {
	tracer := otel.Tracer("accounts-db")
	ctx, span := tracer.Start(ctx, "DB.ConnectAndMigrate")
	defer span.End()

	// 1) open
	dbConn, err := NewDB(dsn)
	if err != nil {
		return nil, err
	}

	// 2) ping
	if err := dbConn.PingContext(ctx); err != nil {
		dbConn.Close()
		return nil, err
	}

	// 3) run migrations
	goose.SetDialect("postgres")
	if err := goose.Up(dbConn, "db/accounts/migration"); err != nil {
		dbConn.Close()
		return nil, err
	}

	return dbConn, nil
}
