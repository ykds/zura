package comet

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/ykds/zura/internal/comet/codec"
	"github.com/ykds/zura/internal/middleware"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/proto/comet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

func NewGrpcServer(srv *Server) *grpc.Server {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(),
			logging.UnaryServerInterceptor(middleware.InterceptorLogger(log.GetGlobalLogger()))),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     time.Minute,
			MaxConnectionAge:      5 * time.Minute,
			MaxConnectionAgeGrace: 3 * time.Second,
			Time:                  10 * time.Second,
			Timeout:               3 * time.Second,
		}))
	comet.RegisterCometServer(server, &GrpcServer{srv: srv})
	listen, err := net.Listen("tcp", "9001")
	if err != nil {
		panic(err)
	}
	go func() {
		err := server.Serve(listen)
		if err != nil {
			log.GetGlobalLogger().Fatalf("comet grpc server exit, error: %+v", err)
		}
	}()
	return server
}

var _ comet.CometServer = &GrpcServer{}

type GrpcServer struct {
	comet.UnimplementedCometServer
	srv *Server
}

func (g *GrpcServer) PushMsg(ctx context.Context, request *comet.PushMsgRequest) (*comet.PushMsgResponse, error) {
	conn, ok := g.srv.onlineUsers[request.Mess.ToUserId]
	if !ok {
		return nil, errors.New(codec.UserIsOffline)
	}
	select {
	case conn.wch <- request.Mess:
	default:
		return nil, errors.New(codec.MessageIsFull)
	}
	return &comet.PushMsgResponse{}, nil
}
