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
	"github.com/ykds/zura/pkg/queue"
	"github.com/ykds/zura/pkg/response"
	"github.com/ykds/zura/proto/logic"
	"github.com/ykds/zura/proto/protocol"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	id          int32
	m           sync.RWMutex
	cfg         *Config
	logicClient logic.LogicClient
	httpServer  http.Server
	onlineUsers map[int64]*Conn
	tq          *queue.TimeQueue
}

func NewServer(c *Config) *Server {
	var engine *gin.Engine
	if c.Server.Debug {
		engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
		engine.Use(gin.RecoveryWithWriter(log.GetGlobalLogger()))
	}
	s := &Server{
		id:  c.Server.ID,
		cfg: c,
		httpServer: http.Server{
			Addr:    ":" + c.HttpServer.Port,
			Handler: engine,
		},
		onlineUsers: make(map[int64]*Conn),
		tq:          queue.NewTimeoutQueue(500, 10240),
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
	go func() {
		s.tq.Run(context.Background())
	}()
	return s
}

func (s *Server) handleWebsocket(c *gin.Context) {
	conn, err := Upgrade(c.Writer, c.Request)
	if err != nil {
		log.Errorf("upgrade websocket failed, err: %v", err)
		return
	}
	err = s.online(c.GetHeader("token"), conn)
	if err != nil {
		log.Errorf("user connect failed. err:%v", err)
		_ = conn.Close()
		return
	}
	go s.Recv(conn)
	go s.Write(conn)
	//if !cfg.Server.Debug {
	//	go s.CheckHeartbeat(conn)
	//}
}

func (s *Server) online(token string, conn *Conn) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()
	resp, err := s.logicClient.Connect(ctx, &logic.ConnectionRequest{
		Token:    token,
		ServerId: s.id,
	})
	if err != nil {
		return err
	}
	conn.UserId = resp.UserId
	s.m.Lock()
	s.onlineUsers[resp.UserId] = conn
	log.Debugf("User[%d] online.", resp.UserId)
	s.m.Unlock()
	return nil
}

func (s *Server) offline(userId int64) error {
	s.m.Lock()
	if _, ok := s.onlineUsers[userId]; !ok {
		s.m.Unlock()
		return nil
	}
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

func (s *Server) CloseConn(conn *Conn) error {
	if conn.isClosed.Load() {
		return nil
	}
	_ = conn.CloseConn()
	_ = s.offline(conn.UserId)
	return nil
}

type Request struct {
	Op   int32           `json:"op"`
	Data json.RawMessage `json:"data"`
}

type Response struct {
	Op   int32       `json:"op"`
	Data interface{} `json:"data"`
}

type AckMessageRequest struct {
	MsgId int64 `json:"msg_id"`
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
			_ = s.CloseConn(conn)
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
				conn.wch <- &protocol.Protocol{
					Op:   protocol.OpErr,
					Body: reply,
				}
				continue
			}
			switch req.Op {
			case protocol.OpAck:
				ackreq := AckMessageRequest{}
				err := json.Unmarshal(req.Data, &ackreq)
				if err != nil {
					log.Errorf("User[%d] ack mess request error, raw content: %s, err: %+v", conn.UserId, string(req.Data), err)
				}
				s.tq.Finish(strconv.FormatInt(ackreq.MsgId, 10))
			case protocol.OpHeartbeat:
				result, err := s.heartbeat(conn.UserId)
				resp := response.GetResponse(err, result)
				reply, _ := json.Marshal(resp)
				_ = conn.WriteMessage(websocket.TextMessage, reply)
				if err != nil {
					_ = s.CloseConn(conn)
					return
				}
				conn.hbticker.Reset(time.Duration(s.cfg.Session.HeartbeatInterval) * time.Second)
			}
		case websocket.CloseMessage:
			_ = s.CloseConn(conn)
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
		case proto := <-conn.wch:
			if proto.Op == protocol.OpNewMsg || proto.Op == protocol.OpNewApplication || proto.Op == protocol.OpApplicationHandlerResult {
				msg := protocol.Message{}
				_ = json.Unmarshal(proto.Body, &msg)
				id := strconv.FormatInt(msg.Id, 10)
				s.tq.Push(id, func() {
					maxTry := 3
					i := 0
					for ; i < maxTry; i++ {
						if s.tq.IsFinished(id) {
							return
						}
						err := conn.WriteMessage(websocket.TextMessage, proto.Body)
						if err != nil {
							log.Errorf("Write message to User[%d] failed, err: %+v", conn.UserId, err)
						}
						time.Sleep(time.Millisecond * 500)
					}
					log.Errorf("User[%d] Retry max send msg, disconnect", conn.UserId)
					_ = s.CloseConn(conn)
				})
			}
			err := conn.WriteMessage(websocket.TextMessage, proto.Body)
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

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
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
		Op: protocol.OpHeartbeatReply,
	}, nil
}
