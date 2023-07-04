package comet

import (
	"github.com/gorilla/websocket"
	"github.com/ykds/zura/proto/protocol"
	"net/http"
	"sync/atomic"
	"time"
)

var upgrader = websocket.Upgrader{}

func Upgrade(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	conn := &Conn{
		Conn:     c,
		wch:      make(chan *protocol.Protocol, 1024),
		close:    make(chan struct{}),
		hbticker: new(time.Ticker),
	}
	conn.isClosed.Store(false)
	return conn, nil
}

type Conn struct {
	*websocket.Conn
	UserId   int64
	wch      chan *protocol.Protocol
	close    chan struct{}
	hbticker *time.Ticker
	isClosed atomic.Bool
}

func (c *Conn) CloseConn() error {
	close(c.close)
	return c.Close()
}
