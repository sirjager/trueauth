package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

// redisCache is a wrapper around redis.Client
type redisCache struct {
	client *redis.Client
	logr   zerolog.Logger
}

var (
	ErrNoRecord  = redis.Nil
	ErrUnMarshal = errors.New("failed to unmarshal data")
)

// NewRedisClient creates a new RedisClient
func NewCacheRedis(client *redis.Client, logr zerolog.Logger) Cache {
	return &redisCache{client, logr}
}

func (r *redisCache) Flush(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}

func (r *redisCache) Set(
	ctx context.Context,
	key string,
	value interface{},
	expiration ...time.Duration,
) error {
	expirationTime := time.Duration(0)
	if len(expiration) > 0 {
		expirationTime = expiration[0]
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, key, data, expirationTime).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) Delete(ctx context.Context, keys ...string) error {
	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisCache) DeleteWithPrefix(ctx context.Context, prefix string) error {
	var cursor uint64 = 0
	var keysToDelete []string
	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, prefix+"*", -1).Result()
		if err != nil {
			return err
		}
		keysToDelete = append(keysToDelete, keys...)
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return r.Delete(ctx, keysToDelete...)
}

func (r *redisCache) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (r *redisCache) GetKeysWithPrefix(ctx context.Context, prefix string) ([]string, error) {
	pattern := prefix + "*"
	return r.GetKeys(ctx, pattern)
}

func (r *redisCache) Get(ctx context.Context, key string, value interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			// can do somethig here, or return different error
			return err // key does not exist
		}
		return err
	}

	if err = json.Unmarshal(data, &value); err != nil {
		return ErrUnMarshal
	}
	return nil
}

func (r *redisCache) GetWithPrefix(ctx context.Context, prefix string, values any) error {
	keys, err := r.GetKeysWithPrefix(ctx, prefix)
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			// can do somethig here, or return different error
			return err // key does not exist
		}
		return err
	}
	if len(keys) == 0 {
		return ErrNoRecord
	}

	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Bytes()
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, values)
		if err != nil {
			return ErrUnMarshal
		}
	}
	return nil
}
