package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

// RedisClient is a wrapper around redis.Client
type RedisClient struct {
	client *redis.Client
	logr   zerolog.Logger
}

// NewRedisClient creates a new RedisClient
func NewRedisClient(address string, logr zerolog.Logger) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{Addr: address})
	return &RedisClient{client, logr}, nil
}

// Ping checks if the Redis server is available
func (r *RedisClient) Ping(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()
	if err != nil {
		r.logr.Error().Err(err).Msg("redis.ping.error")
		return err
	}
	return nil
}

// Close closes the Redis connection
func (r *RedisClient) Close() {
	r.client.Close()
}

// Set sets the value of a key in Redis
func (r *RedisClient) Set(
	ctx context.Context,
	key string,
	data interface{},
	expiration time.Duration,
) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		r.logr.Error().Err(err).Msgf("redis.marshal.error: %v", data)
		return err
	}
	err = r.client.Set(ctx, key, jsonData, expiration).Err()
	if err != nil {
		r.logr.Error().Err(err).Msgf("redis.key: %s", key)
		return err
	}
	return nil
}

// Get gets the value of a key in Redis
func (r *RedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			r.logr.Error().Err(err).Msgf("redis.key: %s", key)
		}
		// no record
		return nil, err
	}
	return []byte(result), nil
}

// Del deletes a key in Redis
func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		r.logr.Error().Err(err).Msgf("redis.marshal.error: %s", keys)
		return err
	}
	return nil
}

// Transaction executes a transaction
func (r *RedisClient) Transaction(
	ctx context.Context,
	commands []redis.Cmder,
) ([]redis.Cmder, error) {
	tx := r.client.TxPipeline()
	for _, cmd := range commands {
		if err := tx.Process(ctx, cmd); err != nil {
			r.logr.Error().Err(err).Msg("redis.transaction.error")
			return nil, err
		}
	}
	_, err := tx.Exec(ctx)
	if err != nil {
		r.logr.Error().Err(err).Msg("redis.transaction.error")
		return nil, err
	}
	return tx.Exec(ctx)
}

// GetKeysWithPrefix returns all keys with a given prefix
func (r *RedisClient) GetKeysWithPrefix(ctx context.Context, prefix string) ([]string, error) {
	var cursor uint64
	var keys []string
	var err error

	for {
		// SCAN command to retrieve keys matching the prefix
		var scanKeys []string
		scanKeys, cursor, err = r.client.Scan(ctx, cursor, prefix+"*", 10).Result()
		if err != nil {
			r.logr.Error().Err(err).Msgf("Failed to retrieve keys with prefix: %s", prefix)
			return nil, err
		}

		// Append the matching keys to the result
		keys = append(keys, scanKeys...)

		// Break the loop if the cursor is 0, indicating the end of the keys
		if cursor == 0 {
			break
		}
	}

	return keys, nil
}
