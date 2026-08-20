package storage

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	if dsn == "" {
		return nil, fmt.Errorf("database DSN is empty")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	if err := createTable(context.Background(), db); err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}
	return &Storage{db: db}, nil
}

func createTable(ctx context.Context, db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS urls (
            id SERIAL PRIMARY KEY,
            short_id TEXT UNIQUE NOT NULL,
            original_url TEXT NOT NULL
        )
    `
	_, err := db.ExecContext(ctx, query)
	return err
}

func (s *Storage) Save(ctx context.Context, url string) (string, error) {
	shortID := generateShortID()
	query := `INSERT INTO urls (short_id, original_url) VALUES ($1, $2)`

	_, err := s.db.ExecContext(ctx, query, shortID, url)
	if err != nil {
		return "", err
	}

	return shortID, nil
}

func (s *Storage) Load(ctx context.Context, shortID string) (string, error) {
	var originalURL string
	query := `SELECT original_url FROM urls WHERE short_id = $1`

	err := s.db.QueryRowContext(ctx, query, shortID).Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, err
}

func (s *Storage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Clear(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "TRUNCATE TABLE urls")
	return err
}

func generateShortID() string {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
