package cache

import (
	"context"
	"time"
)

type Cache interface {
	GetKeys(ctx context.Context, pattern string) ([]string, error)
	GetKeysWithPrefix(ctx context.Context, prefix string) ([]string, error)
	Get(ctx context.Context, key string, value interface{}) error
	GetWithPrefix(ctx context.Context, prefix string, values any) error
	Set(ctx context.Context, key string, value interface{}, expiration ...time.Duration) error
	Delete(ctx context.Context, keys ...string) error
	DeleteWithPrefix(ctx context.Context, prefix string) error
	Flush(ctx context.Context) error
}
