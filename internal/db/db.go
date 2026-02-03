package db

import (
	"context"
	"database/sql"
	"time"
)

// New initializes and returns a new database connection pool.
func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime time.Duration) (*sql.DB, error) {

	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	db.SetConnMaxIdleTime(maxIdleTime)

	// Setting a timeout for the ping operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Pinging the database to ensure the connection is valid
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
