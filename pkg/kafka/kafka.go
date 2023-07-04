package kafka

import (
	"github.com/ykds/zura/pkg/log"
	"sync"
)

type Option func(*Kafka)

func WithLogger(l log.Logger) Option {
	return func(k *Kafka) {
		k.l = l
	}
}

type Kafka struct {
	c         Config
	l         log.Logger
	producers []*Producer
	consumers []*Customer
	m         sync.Mutex
}

func NewKafka(c Config, opts ...Option) *Kafka {
	k := &Kafka{
		c:         c,
		producers: make([]*Producer, 0),
		consumers: make([]*Customer, 0),
	}
	for _, o := range opts {
		o(k)
	}
	return k
}

func (k *Kafka) Close() error {
	k.m.Lock()
	defer k.m.Unlock()
	for _, p := range k.producers {
		_ = p.Close()
	}
	for _, c := range k.consumers {
		_ = c.Close()
	}
	return nil
}
