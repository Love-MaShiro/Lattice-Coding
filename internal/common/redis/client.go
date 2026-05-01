package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"lattice-coding/internal/common/config"
)

type Client struct {
	*goredis.Client
}

func NewClient(cfg *config.RedisConfig) (*Client, error) {
	options := &goredis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		MaxRetries:   cfg.MaxRetries,
	}

	if cfg.DialTimeout != "" {
		if dialTimeout, err := time.ParseDuration(cfg.DialTimeout); err == nil {
			options.DialTimeout = dialTimeout
		}
	}

	if cfg.PoolTimeout != "" {
		if poolTimeout, err := time.ParseDuration(cfg.PoolTimeout); err == nil {
			options.PoolTimeout = poolTimeout
		}
	}

	client := goredis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Client{client}, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}

func (c *Client) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}
