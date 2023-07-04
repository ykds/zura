package server

import (
	"context"
	"github.com/ykds/zura/internal/logic/config"
	"github.com/ykds/zura/internal/logic/services"
	"github.com/ykds/zura/internal/middleware"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/pkg/token"
	"github.com/ykds/zura/proto/logic"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func NewGrpcServer(c config.GrpcServerConfig, service services.Service) *grpc.Server {
	srv := grpc.NewServer(
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
		}),
	)
	logic.RegisterLogicServer(srv, &LogicGrpcServer{srv: service})
	listen, err := net.Listen("tcp", ":"+c.Port)
	if err != nil {
		panic(err)
	}
	go func() {
		err := srv.Serve(listen)
		if err != nil {
			log.Fatalf("logic grpc server exit, error: %+v", err)
		}
	}()
	return srv
}

var _ logic.LogicServer = &LogicGrpcServer{}

type LogicGrpcServer struct {
	logic.UnimplementedLogicServer
	srv services.Service
}

func (l LogicGrpcServer) Connect(ctx context.Context, request *logic.ConnectionRequest) (*logic.ConnectionResponse, error) {
	userId, err := token.VerifyToken(request.Token)
	if err != nil {
		return nil, err
	}
	err = l.srv.UserService.Connect(ctx, userId, request.ServerId)
	if err != nil {
		return &logic.ConnectionResponse{}, err
	}
	return &logic.ConnectionResponse{
		UserId: userId,
	}, nil
}

func (l LogicGrpcServer) Disconnect(ctx context.Context, request *logic.DisconnectRequest) (*logic.DisconnectResponse, error) {
	err := l.srv.UserService.DisConnect(ctx, request.UserId)
	return &logic.DisconnectResponse{}, err
}

func (l LogicGrpcServer) HeartBeat(ctx context.Context, request *logic.HeartBeatRequest) (*logic.HeartBeatResponse, error) {
	err := l.srv.UserService.HeartBeat(ctx, request.UserId)
	return &logic.HeartBeatResponse{}, err
}

func (l LogicGrpcServer) mustEmbedUnimplementedLogicServer() {}
