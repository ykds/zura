package comet

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ykds/zura/internal/comet/codec"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/pkg/pprof"
	"github.com/ykds/zura/pkg/response"
	"github.com/ykds/zura/proto/comet"
	"github.com/ykds/zura/proto/logic"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	m           sync.RWMutex
	cfg         *Config
	logicClient logic.LogicClient
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
		cfg: c,
		httpServer: http.Server{
			Addr:    ":" + c.HttpServer.Port,
			Handler: engine,
		},
		onlineUsers: make(map[int64]*Conn),
	}
	pprof.RouteRegister(engine)
	engine.GET("/ws", s.handleWebsocket)

	conn, err := grpc.DialContext(context.Background(), c.Logic.Host+":"+c.Logic.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		panic(err)
	}
	s.logicClient = logic.NewLogicClient(conn)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("comet server exit, error: %v", err)
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
	auth, err := s.logicClient.Auth(context.Background(), &logic.AuthRequest{
		Token: c.GetHeader("token"),
	})
	if err != nil {
		return
	}
	userId := auth.UserId
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
	if !cfg.Debug {
		go s.CheckHeartbeat(conn)
	}
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
	log.Debugf("User[%d] online.", userId)
	s.m.Unlock()
	return nil
}

func (s *Server) offline(userId int64) error {
	s.m.Lock()
	delete(s.onlineUsers, userId)
	log.Debugf("User[%d] offline.", userId)
	s.m.Unlock()

	if _, err := s.logicClient.Disconnect(context.Background(), &logic.DisconnectRequest{
		UserId: userId,
	}); err != nil {
		return err
	}
	return nil
}

type Request struct {
	Op   comet.Op        `json:"op"`
	Data json.RawMessage `json:"data"`
}

type Response struct {
	Op   comet.Op    `json:"op"`
	Data interface{} `json:"data"`
}

func (s *Server) Recv(conn *Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("User[%d] Recv panic, error: %+v", conn.UserId, err)
		}
	}()
	for {
		message, content, err := conn.ReadMessage()
		if err != nil {
			_ = conn.CloseConn()
			_ = s.offline(conn.UserId)
			log.Errorf("User[%d] err closed connection, error: %+v", conn.UserId, err)
			return
		}

		switch message {
		case websocket.TextMessage:
			req := Request{}
			err = json.Unmarshal(content, &req)
			if err != nil {
				resp := response.GetResponse(errors.Wrap(errors.New(errors.ParameterErrorStatus), err.Error()), nil)
				reply, _ := json.Marshal(resp)
				conn.wch <- reply
				continue
			}
			switch req.Op {
			case comet.Op_SynNewMsg:
				result, err := s.syncMessage(conn.UserId, req.Data)
				var resp response.Resp
				if err != nil {
					resp = response.GetResponse(errors.Wrap(errors.New(codec.SyncNewMessageFailedCode), err.Error()), nil)
				} else {
					resp = response.GetResponse(nil, result)
				}
				reply, _ := json.Marshal(resp)
				conn.wch <- reply
			case comet.Op_Heartbeat:
				result, err := s.heartbeat(conn.UserId)
				resp := response.GetResponse(err, result)
				reply, _ := json.Marshal(resp)
				_ = conn.WriteMessage(websocket.TextMessage, reply)
				if err != nil {
					_ = conn.CloseConn()
					_ = s.offline(conn.UserId)
					return
				}
				conn.hbticker.Reset(time.Duration(s.cfg.Session.HeartbeatInterval) * time.Second)
			}
		case websocket.CloseMessage:
			_ = conn.CloseConn()
			_ = s.offline(conn.UserId)
			log.Debugf("User[%d] initiative closed connection", conn.UserId)
			return
		case websocket.PingMessage:
		}
	}
}

func (s *Server) Write(conn *Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("User[%d] Write panic, error: %+v", conn.UserId, err)
		}
	}()
	for {
		select {
		case message := <-conn.wch:
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Errorf("Write message to User[%d] failed, err: %+v", conn.UserId, err)
				continue
			}
		case <-conn.close:
			err := conn.WriteMessage(websocket.CloseMessage, nil)
			if err != nil {
				log.Errorf("User[%d] connection closed by server", conn.UserId)
			}
			return
		}
	}
}

func (s *Server) CheckHeartbeat(conn *Conn) {
	conn.hbticker = time.NewTicker(time.Duration(s.cfg.Session.HeartbeatInterval) * time.Second)
	for {
		select {
		case <-conn.hbticker.C:
			_ = conn.Close()
			log.Debugf("User[%d] heartbeat timeout", conn.UserId)
		case <-conn.close:
			return
		}
	}
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}

type SyncMessageReq struct {
	SessionId int64 `json:"session_id"`
	Timestamp int64 `json:"timestamp"`
}

func (s *Server) syncMessage(userId int64, body []byte) (*Response, error) {
	req := SyncMessageReq{}
	err := json.Unmarshal(body, &req)
	if err != nil {
		log.Errorf("User[%d] sync mess request error, raw content: %s, err: %+v", userId, string(body), err)
		return nil, err
	}
	listNewMessage, err := s.logicClient.ListNewMessage(context.Background(), &logic.ListNewMessageRequest{
		UserId:    userId,
		SessionId: req.SessionId,
		Timestamp: req.Timestamp,
	})
	if err != nil {
		log.Errorf("User[%d] list new message error, err: %+v", userId, err)
		return nil, err
	}
	return &Response{
		Op:   comet.Op_NewMsgReply,
		Data: listNewMessage.Data,
	}, nil
}

func (s *Server) heartbeat(userId int64) (*Response, error) {
	i := 0
	for ; i < 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_, err := s.logicClient.HeartBeat(ctx, &logic.HeartBeatRequest{UserId: userId})
		if err != nil {
			cancel()
			continue
		}
		cancel()
		break
	}
	if i == 3 {
		return nil, errors.New(codec.HeartBeatFailedCode)
	}
	return &Response{
		Op: comet.Op_HeartbeatReply,
	}, nil
}
