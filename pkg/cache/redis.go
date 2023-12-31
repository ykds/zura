package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	client *redis.Client
	c      Config
	ctx    context.Context
}

func (r Redis) MGet(ctx context.Context, key ...string) ([]string, error) {
	result, err := r.client.MGet(ctx, key...).Result()
	if err != nil {
		return nil, err
	}
	data := make([]string, 0, len(result))
	for _, v := range result {
		if v != nil {
			data = append(data, v.(string))
		} else {
			data = append(data, "")
		}
	}
	return data, nil
}

func (r Redis) LRem(ctx context.Context, key string, start, end int64) error {
	return r.client.LRem(ctx, key, start, end).Err()
}

func (r Redis) LRange(ctx context.Context, key string, start, end int64) ([]string, error) {
	return r.client.LRange(ctx, key, start, end).Result()
}

func (r Redis) LPush(ctx context.Context, key string, value ...interface{}) error {
	return r.client.RPush(ctx, key, value...).Err()
}

func (r Redis) Set(ctx context.Context, key string, value interface{}, ex time.Duration) error {
	return r.client.Set(ctx, key, value, ex).Err()
}

func (r Redis) Get(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			err = NotFoundErr
		}
		return "", err
	}
	return result, nil
}

func (r Redis) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r Redis) Expire(ctx context.Context, key string, ex time.Duration) error {
	return r.client.Expire(ctx, key, ex).Err()
}

var _ Cache = (*Redis)(nil)

func NewRedis(c Config) Cache {
	r := &Redis{
		c: c,
	}
	r.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Host, c.Port),
		Username: c.Username,
		Password: c.Password,
		DB:       c.DB,
	})
	return r
}

func (r Redis) Close() error {
	return r.client.Close()
}
