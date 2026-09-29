package pkg

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func NewRedis(ctx context.Context) (*redis.Client, error) {
	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		return nil, fmt.Errorf("parse redis port: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			os.Getenv("REDIS_HOST"),
			port,
		),
	})

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}
