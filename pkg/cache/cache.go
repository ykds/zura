package cache

import (
	"context"
	"github.com/pkg/errors"
	"sync/atomic"
	"time"
)

var globalCache *atomic.Value

func SetGlobalCache(cache Cache) {
	globalCache = &atomic.Value{}
	globalCache.Store(cache)
}

func GetGlobalCache() Cache {
	if globalCache == nil || globalCache.Load() == nil {
		panic("not initialized")
	}
	return globalCache.Load().(Cache)
}

var NotFoundErr = errors.New("key not found")

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ex time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	MGet(ctx context.Context, key ...string) ([]string, error)
	Del(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, ex time.Duration) error
	LPush(ctx context.Context, key string, value ...interface{}) error
	LRange(ctx context.Context, key string, start, end int64) ([]string, error)
	LRem(ctx context.Context, key string, start, end int64) error
	//HGet(ctx context.Context, key string, field string) (string, error)
	//HMGet(ctx context.Context, key string, field ...string) ([]string, error)
	//HGetAll(ctx context.Context, key string) ([]string, error)
	//HSet(ctx context.Context, key string, kvs ...interface{}) error
}
