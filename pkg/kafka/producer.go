package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	w *kafka.Writer
}

func (p *Producer) WriteMessage(ctx context.Context, key string, value []byte) error {
	return p.w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
}

type ProducerOption func(writer *kafka.Writer)

func WithCustomBalancer(balancer kafka.Balancer) ProducerOption {
	return func(writer *kafka.Writer) {
		writer.Balancer = balancer
	}
}

func (k *Kafka) NewProducer(topic string, opts ...ProducerOption) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(k.c.Brokers...),
		Topic:        topic,
		RequiredAcks: kafka.RequireOne,
	}
	for _, o := range opts {
		o(w)
	}
	return &Producer{w: w}
}
