package ports

import "context"

type Service interface {
	Create(ctx context.Context, originalURL string) (string, error)
	CreateBatch(ctx context.Context, originalURLs []string) ([]string, error)
	Get(ctx context.Context, shortID string) (string, error)
}
