package cache

import (
	"context"
	"github.com/patrickmn/go-cache"
	"time"
)

type Memory struct {
	client *cache.Cache
}

func (m Memory) LRem(ctx context.Context, key string, start, end int64) error {
	v, ex, ok := m.client.GetWithExpiration(key)
	if !ok {
		return NotFoundErr
	}
	l := v.([]string)
	if end < 0 {
		end = int64(len(l)) + end + 1
	}
	l = l[start:end]
	_ = m.client.Replace(key, l, time.Until(ex))
	return nil
}

// TODO L 兼容 interface 各种类型
func (m Memory) LRange(ctx context.Context, key string, start, end int64) ([]string, error) {
	v, ok := m.client.Get(key)
	if !ok {
		return nil, NotFoundErr
	}
	return v.([]string)[start:end], nil
}

func (m Memory) LPush(ctx context.Context, key string, value ...interface{}) error {
	v, ex, ok := m.client.GetWithExpiration(key)
	if !ok {
		l := make([]interface{}, 0)
		v = append(l, value...)
		_ = m.client.Replace(key, l, time.Until(ex))
	} else {
		l := v.([]interface{})
		l = append(l, value...)
		_ = m.client.Replace(key, l, time.Until(ex))
	}
	return nil
}

func (m Memory) Set(ctx context.Context, key string, value interface{}, ex time.Duration) error {
	m.client.Set(key, value, ex)
	return nil
}

func (m Memory) Get(ctx context.Context, key string) (interface{}, error) {
	v, ok := m.client.Get(key)
	if !ok {
		return nil, NotFoundErr
	}
	return v, nil
}

func (m Memory) Del(ctx context.Context, key string) error {
	m.client.Delete(key)
	return nil
}

func (m Memory) Expire(ctx context.Context, key string, ex time.Duration) error {
	v, ok := m.client.Get(key)
	if ok {
		m.client.Set(key, v, ex)
	}
	return nil
}

func NewMemoryCache() Cache {
	globalMemCache = new(Memory)
	globalMemCache.client = cache.New(30*time.Second, time.Minute)
	return globalMemCache
}

var globalMemCache *Memory

func GetGlobalMemCache() Cache {
	if globalMemCache == nil {
		panic("not initialized yes")
	}
	return globalMemCache
}
