package redis

import (
	"context"
	"encoding/json"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"lattice-coding/internal/common/errors"
)

func Get[T any](ctx context.Context, client *Client, key string) (*T, error) {
	data, err := client.Get(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, errors.NotFoundErr("cache not found")
		}
		return nil, errors.CacheErr("get cache failed: " + err.Error())
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, errors.CacheErr("unmarshal cache failed: " + err.Error())
	}

	return &result, nil
}

func Set[T any](ctx context.Context, client *Client, key string, value T, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errors.CacheErr("marshal value failed: " + err.Error())
	}

	if err := client.Set(ctx, key, data, ttl).Err(); err != nil {
		return errors.CacheErr("set cache failed: " + err.Error())
	}

	return nil
}

func SetNX[T any](ctx context.Context, client *Client, key string, value T, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, errors.CacheErr("marshal value failed: " + err.Error())
	}

	result, err := client.SetNX(ctx, key, data, ttl).Result()
	if err != nil {
		return false, errors.CacheErr("setnx failed: " + err.Error())
	}

	return result, nil
}

func Del(ctx context.Context, client *Client, keys ...string) error {
	if err := client.Del(ctx, keys...).Err(); err != nil {
		return errors.CacheErr("delete cache failed: " + err.Error())
	}
	return nil
}

func Exists(ctx context.Context, client *Client, key string) (bool, error) {
	result, err := client.Exists(ctx, key).Result()
	if err != nil {
		return false, errors.CacheErr("check exists failed: " + err.Error())
	}
	return result > 0, nil
}

func Expire(ctx context.Context, client *Client, key string, ttl time.Duration) error {
	if err := client.Expire(ctx, key, ttl).Err(); err != nil {
		return errors.CacheErr("set expire failed: " + err.Error())
	}
	return nil
}

func TTL(ctx context.Context, client *Client, key string) (time.Duration, error) {
	result, err := client.TTL(ctx, key).Result()
	if err != nil {
		return 0, errors.CacheErr("get ttl failed: " + err.Error())
	}
	return result, nil
}

func Incr(ctx context.Context, client *Client, key string) (int64, error) {
	result, err := client.Incr(ctx, key).Result()
	if err != nil {
		return 0, errors.CacheErr("incr failed: " + err.Error())
	}
	return result, nil
}

func Decr(ctx context.Context, client *Client, key string) (int64, error) {
	result, err := client.Decr(ctx, key).Result()
	if err != nil {
		return 0, errors.CacheErr("decr failed: " + err.Error())
	}
	return result, nil
}

func HSet[T any](ctx context.Context, client *Client, key string, field string, value T) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errors.CacheErr("marshal value failed: " + err.Error())
	}

	if err := client.HSet(ctx, key, field, data).Err(); err != nil {
		return errors.CacheErr("hset failed: " + err.Error())
	}
	return nil
}

func HGet[T any](ctx context.Context, client *Client, key string, field string) (*T, error) {
	data, err := client.HGet(ctx, key, field).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, errors.NotFoundErr("hash field not found")
		}
		return nil, errors.CacheErr("hget failed: " + err.Error())
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, errors.CacheErr("unmarshal value failed: " + err.Error())
	}
	return &result, nil
}

func HDel(ctx context.Context, client *Client, key string, fields ...string) error {
	if err := client.HDel(ctx, key, fields...).Err(); err != nil {
		return errors.CacheErr("hdel failed: " + err.Error())
	}
	return nil
}

func LPush[T any](ctx context.Context, client *Client, key string, values ...T) error {
	dataList := make([]interface{}, len(values))
	for i, v := range values {
		data, err := json.Marshal(v)
		if err != nil {
			return errors.CacheErr("marshal value failed: " + err.Error())
		}
		dataList[i] = data
	}

	if err := client.LPush(ctx, key, dataList...).Err(); err != nil {
		return errors.CacheErr("lpush failed: " + err.Error())
	}
	return nil
}

func RPop[T any](ctx context.Context, client *Client, key string) (*T, error) {
	data, err := client.RPop(ctx, key).Bytes()
	if err != nil {
		if err == goredis.Nil {
			return nil, errors.NotFoundErr("list is empty")
		}
		return nil, errors.CacheErr("rpop failed: " + err.Error())
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, errors.CacheErr("unmarshal value failed: " + err.Error())
	}
	return &result, nil
}
