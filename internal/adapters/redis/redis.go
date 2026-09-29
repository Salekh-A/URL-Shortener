package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const cacheTTL = time.Hour

type Cache struct {
	client *redis.Client
}

func New(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Cache) Set(ctx context.Context, key string, value string) error {
	return c.client.Set(ctx, key, value, cacheTTL).Err()
}
