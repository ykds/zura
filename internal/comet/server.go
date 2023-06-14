package comet

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/pkg/response"
	"github.com/ykds/zura/proto/logic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"sync"
	"time"
)

var cfg = new(Config)

func GetConfig() *Config {
	return cfg
}

type Config struct {
	Debug bool       `json:"debug" yaml:"debug"`
	Port  string     `json:"port" yaml:"port"`
	Logic Logic      `json:"logic" yaml:"logic"`
	Log   log.Config `json:"log" yaml:"log"`
}

type Logic struct {
	Host string `json:"host" yaml:"host"`
	Port string `json:"port" yaml:"port"`
}

type Server struct {
	logicClient logic.LogicClient
	m           sync.RWMutex
	httpServer  http.Server
	onlineUsers map[int64]*Conn
}

func NewServer(c *Config) *Server {
	var engine *gin.Engine
	if c.Debug {
		engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
		engine.Use(gin.LoggerWithWriter(log.GetGlobalLogger()), gin.RecoveryWithWriter(log.GetGlobalLogger()))
	}
	s := &Server{
		httpServer: http.Server{
			Addr:    ":" + c.Port,
			Handler: engine,
		},
		onlineUsers: make(map[int64]*Conn),
	}
	engine.GET("/ws", s.handleWebsocket)

	conn, err := grpc.DialContext(context.Background(), c.Logic.Host+":"+c.Logic.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	s.logicClient = logic.NewLogicClient(conn)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.GetGlobalLogger().Fatalf("comet server exit, error: %v", err)
			return
		}
	}()
	return s
}

func (s *Server) handleWebsocket(c *gin.Context) {
	var (
		err error
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	userId := c.GetInt64(common.UserIdKey)
	conn, err := Upgrade(c.Writer, c.Request)
	if err != nil {
		return
	}
	err = s.online(userId, conn)
	if err != nil {
		return
	}
	go s.Recv(conn)
	go s.Write(conn)
}

func (s *Server) online(userId int64, conn *Conn) error {
	if _, err := s.logicClient.Connect(context.Background(), &logic.ConnectionRequest{
		UserId: userId,
	}); err != nil {
		return err
	}
	conn.UserId = userId
	s.m.Lock()
	s.onlineUsers[userId] = conn
	s.m.Unlock()
	return nil
}

func (s *Server) offline(userId int64) error {
	s.m.RLock()
	_, ok := s.onlineUsers[userId]
	if !ok {
		s.m.RUnlock()
		return nil
	}
	s.m.RUnlock()
	if _, err := s.logicClient.Disconnect(context.Background(), &logic.DisconnectRequest{
		UserId: userId,
	}); err != nil {
		return err
	}
	s.m.Lock()
	delete(s.onlineUsers, userId)
	s.m.Unlock()
	return nil
}

func (s *Server) Recv(conn *Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.GetGlobalLogger().Errorf("User[%d] Recv panic, error: %+v", conn.UserId, err)
		}
	}()
	for {
		message, _, err := conn.ReadMessage()
		if err != nil {
			_ = conn.CloseConn()
			log.GetGlobalLogger().Errorf("User[%d] err closed connection, error: %+v", conn.UserId, err)
			return
		}
		switch message {
		case websocket.TextMessage:
		case websocket.CloseMessage:
			_ = conn.CloseConn()
			log.GetGlobalLogger().Infof("User[%d] initiative closed connection", conn.UserId)
			return
		case websocket.PingMessage:
			i := 0
			for ; i < 3; i++ {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				_, err := s.logicClient.HeartBeat(ctx, &logic.HeartBeatRequest{UserId: conn.UserId})
				if err != nil {
					cancel()
					continue
				}
				cancel()
				_ = conn.WriteMessage(websocket.PongMessage, nil)
				break
			}
			if i == 3 {
				_ = conn.CloseConn()
				log.GetGlobalLogger().Infof("User[%d] heartbeat failed, err: %+V", conn.UserId, err)
				return
			}
		}
	}
}

func (s *Server) Write(conn *Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.GetGlobalLogger().Errorf("User[%d] Write panic, error: %+v", conn.UserId, err)
		}
	}()
	for {
		select {
		case message := <-conn.wch:
			data, _ := json.Marshal(message)
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.GetGlobalLogger().Errorf("Write message to User[%d] failed, err: %+v", conn.UserId, err)
			}
		case <-conn.close:
			err := conn.WriteMessage(websocket.CloseMessage, nil)
			if err != nil {
				log.GetGlobalLogger().Infof("User[%d] connection closed by server", conn.UserId)
				return
			}
		}
	}
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
