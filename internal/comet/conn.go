package comet

import (
	"github.com/gorilla/websocket"
	"net/http"
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
		wch:      make(chan []byte, 10),
		close:    make(chan struct{}),
		hbticker: new(time.Ticker),
	}
	return conn, nil
}

type Conn struct {
	*websocket.Conn
	UserId   int64
	wch      chan []byte
	close    chan struct{}
	hbticker *time.Ticker
}

func (c *Conn) CloseConn() error {
	close(c.close)
	return c.Close()
}
