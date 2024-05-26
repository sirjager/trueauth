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

var ErrNoRecord = redis.Nil

// NewRedisClient creates a new RedisClient
func NewCacheRedis(client *redis.Client, logr zerolog.Logger) Cache {
	return &redisCache{client, logr}
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
		r.logr.Error().Err(err).Msgf("failed to marshal redis cache for key::%s", key)
		return err
	}
	err = r.client.Set(ctx, key, data, expirationTime).Err()
	if err != nil {
		r.logr.Error().Err(err).Msgf("failed to set redis cache for key::%s", key)
		return err
	}
	return nil
}

func (r *redisCache) Get(ctx context.Context, key string, value any) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, ErrNoRecord) {
			r.logr.Error().Err(err).Msgf("no record found for key::%s", key)
			return err
		}
		r.logr.Error().Err(err).Msgf("failed to get redis cache for key::%s", key)
		return err
	}
	return json.Unmarshal(data, value)
}

func (r *redisCache) Delete(ctx context.Context, keys ...string) error {
	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		r.logr.Error().Err(err).Msgf("failed to delete redis cache for key::%v", keys)
		return err
	}
	return nil
}

func (r *redisCache) Flush(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}
