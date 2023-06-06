package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func Upgrade(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return &Conn{
		c:   c,
		wch: make(chan Message, 10),
	}, nil
}

type Conn struct {
	UserId int64
	c      *websocket.Conn
	wch    chan Message
}

type Message struct {
	FromUserId int64  `json:"from_user_id"`
	ToUserId   int64  `json:"to_user_id"`
	Timestamp  int64  `json:"timestamp"`
	Content    string `json:"content"`
}

func (c *Conn) Recv() {
	for {
		message, _, err := c.c.ReadMessage()
		if err != nil {
			_ = c.Close()
			return
		}
		switch message {
		case websocket.TextMessage:
			// 用 chan 处理
			// TODO 业务处理
		case websocket.CloseMessage:
			_ = c.Close()
			return
		case websocket.PingMessage:
			// 用 chan 处理
			// TODO 心跳
		}
	}
}

func (c *Conn) Write() {
	for {
		for m := range c.wch {
			data, _ := json.Marshal(m)
			err := c.c.WriteMessage(websocket.TextMessage, data)
			// TODO log or return?
			if err != nil {
				return
			}
		}
	}
}

func (c *Conn) Close() error {
	return c.c.Close()
}
