package ports

import (
	"context"

	"newproject/internal/domain"
)

type Repository interface {
	Save(ctx context.Context, url domain.URL) error
	SaveBatch(ctx context.Context, urls []domain.URL) error
	GetByShortID(ctx context.Context, shortID string) (*domain.URL, error)
	GetByOriginalURL(ctx context.Context, originalURL string) (*domain.URL, error)
}
