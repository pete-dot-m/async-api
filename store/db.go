package store

import (
	"async-api/config"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

func NewPostgresDB(conf *config.Config) (*sql.DB, error) {
	dsn := conf.DatabaseUrl()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
