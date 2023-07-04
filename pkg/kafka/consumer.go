package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/ykds/zura/pkg/log"
	"os"
	"time"
)

type ConsumerHandler interface {
	Consume(kafka.Message) error
}

type ConsumerOption func(config *kafka.ReaderConfig)

func WithCommitInterval(interval time.Duration) ConsumerOption {
	return func(config *kafka.ReaderConfig) {
		config.CommitInterval = interval
	}
}

func WithPartition(partition int) ConsumerOption {
	return func(config *kafka.ReaderConfig) {
		config.Partition = partition
	}
}

func (k *Kafka) NewConsumer(groupId string, topic []string, consume ConsumerHandler, opts ...ConsumerOption) *Customer {
	c := kafka.ReaderConfig{
		Brokers:     k.c.Brokers,
		GroupID:     groupId,
		GroupTopics: topic,
	}
	for _, o := range opts {
		o(&c)
	}
	con := &Customer{
		l:       k.l,
		r:       kafka.NewReader(c),
		consume: consume,
	}
	k.m.Lock()
	defer k.m.Unlock()
	k.consumers = append(k.consumers, con)
	return con
}

type Customer struct {
	l       log.Logger
	r       *kafka.Reader
	consume ConsumerHandler
}

func (c *Customer) Run(ctx context.Context) {
	for {
		message, err := c.r.FetchMessage(ctx)
		if err != nil {
			if c.l != nil {
				c.l.Error(err)
			}
		}
		err = c.consume.Consume(message)
		if err != nil {
			if c.l != nil {
				c.l.Error(err)
			} else {
				_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf("%+v", err))
			}
		}
		err = c.r.CommitMessages(ctx, message)
		if err != nil {
			if c.l != nil {
				c.l.Error(err)
			}
		}

		if ctx.Err() != nil {
			return
		}
	}
}

func (c *Customer) Close() error {
	return c.r.Close()
}
