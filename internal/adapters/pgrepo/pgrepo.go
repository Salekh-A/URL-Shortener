package pgrepo

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"newproject/internal/domain"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Save(ctx context.Context, url domain.URL) error {
	query := `INSERT INTO urls (short_id, original_url) VALUES ($1, $2)`

	_, err := r.db.Exec(ctx, query, url.ShortID, url.OriginalURL)
	return err
}

func (r *Repository) SaveBatch(ctx context.Context, urls []domain.URL) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}

	for _, url := range urls {
		batch.Queue(
			`INSERT INTO urls (short_id, original_url) VALUES ($1, $2)`,
			url.ShortID,
			url.OriginalURL,
		)
	}

	results := tx.SendBatch(ctx, batch)

	for range urls {
		if _, err := results.Exec(); err != nil {
			results.Close()
			return err
		}
	}

	if err := results.Close(); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetByShortID(ctx context.Context, shortID string) (*domain.URL, error) {
	var url domain.URL

	query := `SELECT short_id, original_url FROM urls WHERE short_id = $1`

	err := r.db.QueryRow(ctx, query, shortID).Scan(
		&url.ShortID,
		&url.OriginalURL,
	)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *Repository) GetByOriginalURL(ctx context.Context, originalURL string) (*domain.URL, error) {
	var url domain.URL

	query := `
		SELECT short_id, original_url
		FROM urls
		WHERE original_url = $1
	`

	err := r.db.QueryRow(ctx, query, originalURL).Scan(
		&url.ShortID,
		&url.OriginalURL,
	)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func (r *Repository) Clear(ctx context.Context) error {
	_, err := r.db.Exec(ctx, "TRUNCATE TABLE urls")
	return err
}
