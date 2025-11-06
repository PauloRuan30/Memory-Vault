package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewPostgresDB(host, port, user, password, dbName string) (*PostgresDB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	// Test the connection
	if err = pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Create tables if they don't exist
	if err = createTables(pool); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return &PostgresDB{Pool: pool}, nil
}

func createTables(pool *pgxpool.Pool) error {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);`

	filesTable := `
	CREATE TABLE IF NOT EXISTS files (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		parent_folder_id INTEGER REFERENCES files(id) ON DELETE CASCADE,
		is_folder BOOLEAN NOT NULL DEFAULT false,
		name VARCHAR(255) NOT NULL,
		s3_path VARCHAR(1024),
		texture_path VARCHAR(1024),
		size_kb INTEGER NOT NULL DEFAULT 0,
		processing_status VARCHAR(20) DEFAULT 'PENDING',
		metadata JSONB,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);`

	_, err := pool.Exec(context.Background(), usersTable)
	if err != nil {
		return err
	}

	_, err = pool.Exec(context.Background(), filesTable)
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresDB) Close() {
	db.Pool.Close()
}
