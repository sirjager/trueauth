package db

import (
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type RedisClient struct {
	Client *redis.Client
	Logr   zerolog.Logger
}

func NewRedisClient(address string, logger zerolog.Logger) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{Addr: address})
	return &RedisClient{Client: client, Logr: logger}, nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}
