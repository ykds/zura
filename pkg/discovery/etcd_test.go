package discovery

import (
	"context"
	"testing"
)

func TestEtcd(t *testing.T) {
	etcd := NewEtcd(Config{
		Urls: []string{"http://localhost:2379"},
	}, "")
	err := etcd.Register(context.Background(), "/zura/comet/test", "127.0.0.1:8000", nil)
	if err != nil {
		panic(err)
	}
}
