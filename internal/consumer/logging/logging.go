package logging

import (
	kafka2 "github.com/segmentio/kafka-go"
	"github.com/ykds/zura/pkg/kafka"
	"net"
)

type loggingConsumer struct {
	logstashConn net.Conn
}

func NewLoggingConsumer() kafka.ConsumerHandler {
	conn, err := net.Dial("tcp", ":4560")
	if err != nil {
		panic(err)
	}
	return &loggingConsumer{logstashConn: conn}
}

func (l *loggingConsumer) Consume(msg kafka2.Message) error {
	_, err := l.logstashConn.Write(msg.Value)
	return err
}
