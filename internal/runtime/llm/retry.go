package llm

import (
	"context"
	"math"
	"math/rand"
	"time"
)

type RetryConfig struct {
	MaxAttempts     int
	InitialInterval time.Duration
	MaxInterval    time.Duration
}

func NewRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:     2,
		InitialInterval: 500 * time.Millisecond,
		MaxInterval:    5 * time.Second,
	}
}

type RetryableFunc func(ctx context.Context) error

func DoWithRetry(ctx context.Context, cfg RetryConfig, fn RetryableFunc) error {
	var err error
	interval := cfg.InitialInterval

	for attempt := 0; attempt <= cfg.MaxAttempts; attempt++ {
		if attempt > 0 && interval > 0 {
			jitter := time.Duration(rand.Int63n(int64(interval / 2)))
			sleep := interval + jitter

			select {
			case <-time.After(sleep):
			case <-ctx.Done():
				return ctx.Err()
			}

			interval = time.Duration(math.Min(float64(interval)*2, float64(cfg.MaxInterval)))
		}

		err = fn(ctx)
		if err == nil {
			return nil
		}

		if !isRetryable(err) {
			return err
		}
	}

	return err
}

func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	return true
}
