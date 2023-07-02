package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"testing"
)

var _ ConsumerHandler = new(consumer)

type consumer struct{}

func (c consumer) Consume(message kafka.Message) error {
	fmt.Println(string(message.Key), string(message.Value))
	return nil
}

func TestKafka(t *testing.T) {
	k := NewKafka(DefaultConfig())
	produce := k.NewProducer("test1")
	err := produce.WriteMessage(context.Background(), "hello", []byte("world"))
	if err != nil {
		panic(err)
	}
	err = produce.WriteMessage(context.Background(), "hello1", []byte("world2"))
	if err != nil {
		panic(err)
	}
	c := k.NewConsumer("grouptest", []string{"test1"}, consumer{})
	c.Run(nil)
}
