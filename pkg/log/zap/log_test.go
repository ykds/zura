package zap

import (
	"github.com/ykds/zura/pkg/log"
	"testing"
)

func TestLogger(t *testing.T) {
	l := NewLogger(log.DefaultConfig(), WithLogstash(), WithDebug(true))
	l.Errorf("test logging")
}
