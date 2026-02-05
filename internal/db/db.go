package db

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// New initializes and returns a new database connection pool.
func setupDatabaseConn(addr string, maxOpenConns, maxIdleConns int, maxIdleTime time.Duration) (*sql.DB, error) {

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

func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime time.Duration) (*sql.DB, error) {

	db, err := setupDatabaseConn(addr, maxOpenConns, maxIdleConns, maxIdleTime)
	if err != nil {
		log.Panic("Error connecting to the database:", err)
	}

	log.Println("Connected to the database successfully")

	return db, nil

}
