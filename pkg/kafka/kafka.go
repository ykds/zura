package kafka

import (
	"github.com/ykds/zura/pkg/log"
)

type Option func(*Kafka)

func WithLogger(l log.Logger) Option {
	return func(k *Kafka) {
		k.l = l
	}
}

type Kafka struct {
	c Config
	l log.Logger
}

func NewKafka(c Config, opts ...Option) *Kafka {
	k := &Kafka{
		c: c,
	}
	for _, o := range opts {
		o(k)
	}
	return k
}
