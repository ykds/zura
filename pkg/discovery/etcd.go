package discovery

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"sync"
)

type Config struct {
	Urls []string `json:"urls"`
}

type EtcdManager struct {
	client      *clientv3.Client
	registerMap map[string]context.CancelFunc
	m           sync.Mutex

	manager  endpoints.Manager
	endpoint string
}

func NewEtcd(c Config, endpoint string) *EtcdManager {
	manager := &EtcdManager{
		registerMap: make(map[string]context.CancelFunc),
	}
	cli, err := clientv3.NewFromURLs(c.Urls)
	if err != nil {
		panic(err)
	}
	manager.client = cli

	mag, err := endpoints.NewManager(cli, endpoint)
	if err != nil {
		panic(err)
	}
	manager.manager = mag
	manager.endpoint = endpoint
	return manager
}

func (e *EtcdManager) Register(ctx context.Context, key, addr string, md map[string]interface{}) error {
	alive, err := e.putKeyWithLease(key, addr, md)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(ctx)
	e.m.Lock()
	e.registerMap[key] = cancel
	e.m.Unlock()
	for {
		select {
		case resp := <-alive:
			if ctx.Err() != nil {
				return ctx.Err()
			}
			// keepalive 失效后返回 nil
			if resp == nil {
				alive, err = e.putKeyWithLease(key, addr, md)
				if err != nil {
					return err
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (e *EtcdManager) UnRegister(key string) {
	e.m.Lock()
	if cancel, ok := e.registerMap[key]; ok {
		cancel()
		delete(e.registerMap, key)
	}
	e.m.Unlock()
}

func (e *EtcdManager) Watch(ctx context.Context, callback func([]*endpoints.Update)) error {
	channel, err := e.manager.NewWatchChannel(ctx)
	if err != nil {
		return err
	}
	for {
		select {
		case resp := <-channel:
			callback(resp)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (e *EtcdManager) putKeyWithLease(key, addr string, md map[string]interface{}) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	lease := clientv3.NewLease(e.client)
	grant, err := lease.Grant(context.Background(), 30)
	if err != nil {
		return nil, err
	}
	err = e.manager.AddEndpoint(context.Background(), e.endpoint+key, endpoints.Endpoint{Addr: addr, Metadata: md}, clientv3.WithLease(grant.ID))
	if err != nil {
		return nil, err
	}
	return lease.KeepAlive(context.Background(), grant.ID)
}

func (e *EtcdManager) Close() error {
	return e.client.Close()
}
