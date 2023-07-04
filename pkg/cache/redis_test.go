package cache

import (
	"context"
	"testing"
)

func TestNewRedis(t *testing.T) {
	redis := NewRedis(DefaultConfig())

	_, _ = redis.MGet(context.Background(), "demo", "test")
}
