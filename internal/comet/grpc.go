package comet

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/ykds/zura/internal/middleware"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/proto/comet"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

func NewGrpcServer(srv *Server) *grpc.Server {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			otelgrpc.UnaryServerInterceptor(),
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
	listen, err := net.Listen("tcp", ":"+srv.cfg.GrpcServer.Port)
	if err != nil {
		panic(err)
	}
	go func() {
		err := server.Serve(listen)
		if err != nil {
			log.Fatalf("comet grpc server exit, error: %+v", err)
		}
	}()
	return server
}

var _ comet.CometServer = &GrpcServer{}

type GrpcServer struct {
	comet.UnimplementedCometServer
	srv *Server
}

func (g *GrpcServer) PushMessage(ctx context.Context, request *comet.PushMsgRequest) (*comet.PushMsgResponse, error) {
	for _, id := range request.ToUserId {
		conn, ok := g.srv.onlineUsers[id]
		if ok {
			conn.wch <- request.Proto
		}
	}
	return &comet.PushMsgResponse{}, nil
}

func (g *GrpcServer) PushNotification(ctx context.Context, request *comet.PushNotificationRequest) (*comet.PushNotificationResponse, error) {
	for _, id := range request.ToUserId {
		conn, ok := g.srv.onlineUsers[id]
		if ok {
			conn.wch <- request.Proto
		}
	}
	return &comet.PushNotificationResponse{}, nil
}
