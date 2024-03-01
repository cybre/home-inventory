package cache

import (
	"context"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
)

type Cache[T any] struct {
	cache      cache.CacheInterface[T]
	expiration time.Duration
}

func New[T any](cache cache.CacheInterface[T], defaultExpiration time.Duration) *Cache[T] {
	return &Cache[T]{
		cache:      cache,
		expiration: defaultExpiration,
	}
}

func (c *Cache[T]) Get(ctx context.Context, key any) (T, error) {
	return c.cache.Get(ctx, key)
}

func (c *Cache[T]) Set(ctx context.Context, key any, object T) error {
	return c.cache.Set(ctx, key, object, store.WithExpiration(c.expiration))
}

func (c *Cache[T]) SetWithExpiration(ctx context.Context, key any, object T, expiration time.Duration) error {
	return c.cache.Set(ctx, key, object, store.WithExpiration(expiration))
}

func (c *Cache[T]) Delete(ctx context.Context, key any) error {
	return c.cache.Delete(ctx, key)
}
