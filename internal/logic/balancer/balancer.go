package balancer

import (
	"encoding/binary"
	"github.com/segmentio/kafka-go"
)

func SessionIdBalance(msg kafka.Message, partitions ...int) (partition int) {
	sessionId := binary.BigEndian.Uint64(msg.Key)
	return int(sessionId % uint64(len(partitions)))
}
