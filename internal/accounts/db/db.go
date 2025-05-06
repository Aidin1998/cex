package db

import (
	"context"
	"database/sql"
	"fmt"

	"cex/pkg/cfg"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

// ConnectAndMigrate connects to the database and runs migrations.
func ConnectAndMigrate(ctx context.Context) (*sql.DB, error) {
	// Read DSN from configuration
	dsn := cfg.Cfg.DB.URL
	if dsn == "" {
		return nil, fmt.Errorf("database DSN is not configured")
	}

	// Open a connection to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify the connection
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set Goose dialect and run migrations
	goose.SetDialect("postgres")
	if err := goose.Up(db, "db/accounts/migration"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}
