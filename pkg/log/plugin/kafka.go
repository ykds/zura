package plugin

import (
	"context"
	"github.com/ykds/zura/pkg/kafka"
	"io"
)

type kafkaWriter struct {
	*kafka.Producer
}

func (k kafkaWriter) Write(p []byte) (n int, err error) {
	_ = k.WriteMessage(context.TODO(), "", p)
	return 0, nil
}

func NewKafkaWriter(producer *kafka.Producer) io.Writer {
	return &kafkaWriter{Producer: producer}
}
