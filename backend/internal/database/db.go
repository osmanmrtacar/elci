package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

// New creates a new database connection
func New(dbPath string) (*DB, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	sqlDB, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{DB: sqlDB}

	// Run migrations
	if err := db.RunMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database connected successfully")
	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// Begin starts a new transaction
func (db *DB) Begin() (*sql.Tx, error) {
	return db.DB.Begin()
}

// Health checks if the database connection is healthy
func (db *DB) Health() error {
	return db.Ping()
}
