package cache

import (
	"context"
	"github.com/pkg/errors"
	"time"
)

var NotFoundErr = errors.New("key not found")

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ex time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	Del(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, ex time.Duration) error
	LPush(ctx context.Context, key string, value ...interface{}) error
	LRange(ctx context.Context, key string, start, end int64) ([]string, error)
	LRem(ctx context.Context, key string, start, end int64) error
}
