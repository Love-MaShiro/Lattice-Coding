package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"lattice-coding/internal/common/errors"
)

type Lock struct {
	client *Client
	key    string
	value  string
	ttl    time.Duration
}

func NewLock(client *Client, key string, value string, ttl time.Duration) *Lock {
	return &Lock{
		client: client,
		key:    key,
		value:  value,
		ttl:    ttl,
	}
}

func (l *Lock) Acquire(ctx context.Context) (bool, error) {
	result, err := l.client.SetNX(ctx, l.key, l.value, l.ttl).Result()
	if err != nil {
		return false, errors.CacheErr("acquire lock failed: " + err.Error())
	}
	return result, nil
}

func (l *Lock) Release(ctx context.Context) error {
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`

	_, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Result()
	if err != nil && err != goredis.Nil {
		return errors.CacheErr("release lock failed: " + err.Error())
	}
	return nil
}

func (l *Lock) Extend(ctx context.Context, ttl time.Duration) error {
	script := `
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("PEXPIRE", KEYS[1], ARGV[2])
		else
			return 0
		end
	`

	_, err := l.client.Eval(ctx, script, []string{l.key}, l.value, ttl.Milliseconds()).Result()
	if err != nil && err != goredis.Nil {
		return errors.CacheErr("extend lock failed: " + err.Error())
	}
	return nil
}

func TryLock(ctx context.Context, client *Client, key string, value string, ttl time.Duration, maxRetries int, retryInterval time.Duration) (*Lock, error) {
	lock := NewLock(client, key, value, ttl)

	for i := 0; i < maxRetries; i++ {
		acquired, err := lock.Acquire(ctx)
		if err != nil {
			return nil, err
		}
		if acquired {
			return lock, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(retryInterval):
		}
	}

	return nil, errors.CacheErr("failed to acquire lock after max retries")
}
