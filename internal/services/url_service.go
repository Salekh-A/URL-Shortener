package services

import (
	"context"

	"newproject/internal/domain"
	"newproject/internal/ports"
)

type URLService struct {
	repository ports.Repository
	cache      ports.Cache
}

func NewURLService(
	repository ports.Repository,
	cache ports.Cache,
) *URLService {
	return &URLService{
		repository: repository,
		cache:      cache,
	}
}

func (s *URLService) Create(ctx context.Context, originalURL string) (string, error) {
	shortID, err := domain.GenerateShortID()
	if err != nil {
		return "", err
	}

	url := domain.URL{
		ShortID:     shortID,
		OriginalURL: originalURL,
	}

	if err := s.repository.Save(ctx, url); err != nil {
		return "", err
	}

	return shortID, nil
}

func (s *URLService) CreateBatch(ctx context.Context, originalURLs []string) ([]string, error) {
	urls := make([]domain.URL, 0, len(originalURLs))
	shortIDs := make([]string, 0, len(originalURLs))

	for _, originalURL := range originalURLs {
		shortID, err := domain.GenerateShortID()
		if err != nil {
			return nil, err
		}

		urls = append(urls, domain.URL{
			ShortID:     shortID,
			OriginalURL: originalURL,
		})

		shortIDs = append(shortIDs, shortID)
	}

	if err := s.repository.SaveBatch(ctx, urls); err != nil {
		return nil, err
	}

	return shortIDs, nil
}

func (s *URLService) Get(ctx context.Context, shortID string) (string, error) {
	cacheKey := "url:" + shortID

	cachedURL, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		return cachedURL, nil
	}

	url, err := s.repository.GetByShortID(ctx, shortID)
	if err != nil {
		return "", err
	}

	_ = s.cache.Set(ctx, cacheKey, url.OriginalURL)

	return url.OriginalURL, nil
}
