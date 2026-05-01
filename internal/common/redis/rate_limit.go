package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"lattice-coding/internal/common/errors"
)

type RateLimiter struct {
	client       *Client
	key          string
	maxRequests  int
	window       time.Duration
}

func NewRateLimiter(client *Client, key string, maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		client:       client,
		key:          key,
		maxRequests:  maxRequests,
		window:       window,
	}
}

func (r *RateLimiter) Allow(ctx context.Context) (bool, error) {
	key := fmt.Sprintf("%s:%d", r.key, time.Now().Unix()/int64(r.window.Seconds()))

	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, errors.CacheErr("rate limit check failed: " + err.Error())
	}

	if count == 1 {
		r.client.Expire(ctx, key, r.window)
	}

	return count <= int64(r.maxRequests), nil
}

func (r *RateLimiter) Remaining(ctx context.Context) (int, error) {
	key := fmt.Sprintf("%s:%d", r.key, time.Now().Unix()/int64(r.window.Seconds()))

	count, err := r.client.Get(ctx, key).Int()
	if err != nil && err != goredis.Nil {
		return 0, errors.CacheErr("rate limit remaining check failed: " + err.Error())
	}

	remaining := r.maxRequests - count
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}

func (r *RateLimiter) Reset(ctx context.Context) error {
	key := fmt.Sprintf("%s:%d", r.key, time.Now().Unix()/int64(r.window.Seconds()))
	return r.client.Del(ctx, key).Err()
}

type SlidingWindowLimiter struct {
	client       *Client
	key          string
	maxRequests  int
	window       time.Duration
}

func NewSlidingWindowLimiter(client *Client, key string, maxRequests int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		client:       client,
		key:          key,
		maxRequests:  maxRequests,
		window:       window,
	}
}

func (l *SlidingWindowLimiter) Allow(ctx context.Context) (bool, error) {
	now := time.Now().UnixMilli()
	windowStart := now - l.window.Milliseconds()

	pipe := l.client.Pipeline()

	pipe.ZRemRangeByScore(ctx, l.key, "0", fmt.Sprintf("%d", windowStart))

	countCmd := pipe.ZCard(ctx, l.key)

	pipe.ZAdd(ctx, l.key, goredis.Z{Score: float64(now), Member: now})

	pipe.Expire(ctx, l.key, l.window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, errors.CacheErr("sliding window check failed: " + err.Error())
	}

	count := countCmd.Val()
	return count < int64(l.maxRequests), nil
}

func (l *SlidingWindowLimiter) Remaining(ctx context.Context) (int, error) {
	now := time.Now().UnixMilli()
	windowStart := now - l.window.Milliseconds()

	count, err := l.client.ZCount(ctx, l.key, fmt.Sprintf("%d", windowStart), fmt.Sprintf("%d", now)).Result()
	if err != nil {
		return 0, errors.CacheErr("sliding window remaining check failed: " + err.Error())
	}

	remaining := l.maxRequests - int(count)
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}
