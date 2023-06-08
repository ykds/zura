package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	*redis.Client
	c   *Config
	ctx context.Context
}

func NewRedis(ctx context.Context, c *Config) *Redis {
	r := &Redis{
		c:   c,
		ctx: ctx,
	}
	r.Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Host, c.Port),
		Username: c.Username,
		Password: c.Password,
		DB:       c.DB,
	})
	go func() {
		select {
		case <-ctx.Done():
			r.Client.Close()
		}
	}()
	return r
}
