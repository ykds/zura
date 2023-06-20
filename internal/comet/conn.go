package comet

import (
	"github.com/gorilla/websocket"
	"github.com/ykds/zura/proto/comet"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func Upgrade(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	conn := &Conn{
		Conn:  c,
		wch:   make(chan *comet.Proto, 10),
		close: make(chan struct{}),
	}
	return conn, nil
}

type Conn struct {
	*websocket.Conn
	UserId int64
	wch    chan *comet.Proto
	close  chan struct{}
}

func (c *Conn) CloseConn() error {
	close(c.close)
	return c.Close()
}
