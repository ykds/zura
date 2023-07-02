package cache

import (
	"context"
	"testing"
)

func TestNewRedis(t *testing.T) {
	redis := NewRedis(DefaultConfig())
	err := redis.Set(context.Background(), "demo", false, 0)
	if err != nil {
		t.Error(err)
	}
	redis.Get(context.Background(), "demo")
}
