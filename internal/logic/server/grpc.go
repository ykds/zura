package server

import (
	"context"
	"github.com/ykds/zura/internal/logic/services"
	"github.com/ykds/zura/internal/logic/services/message"
	"github.com/ykds/zura/internal/middleware"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/pkg/token"
	"github.com/ykds/zura/proto/logic"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
)

type GrpcServerConfig struct {
	Port string `json:"port" yaml:"port"`
}

func DefaultGrpcConfig() GrpcServerConfig {
	return GrpcServerConfig{
		Port: "8001",
	}
}

func NewGrpcServer(c GrpcServerConfig, service services.Service) *grpc.Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
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

func (l LogicGrpcServer) Auth(ctx context.Context, request *logic.AuthRequest) (*logic.AuthResponse, error) {
	userId, err := token.VerifyToken(request.Token)
	if err != nil {
		return &logic.AuthResponse{}, err
	}
	return &logic.AuthResponse{UserId: userId}, nil
}

func (l LogicGrpcServer) ListNewMessage(ctx context.Context, request *logic.ListNewMessageRequest) (*logic.ListNewMessageResponse, error) {
	newMessage, err := l.srv.MessageService.ListNewMessage(request.UserId, message.ListMessageRequest{
		SessionId: request.SessionId,
		Timestamp: request.Timestamp,
	})
	if err != nil {
		return &logic.ListNewMessageResponse{}, err
	}
	data := make([]*logic.MessageItem, 0, len(newMessage))
	for _, item := range newMessage {
		data = append(data, &logic.MessageItem{
			Id:         item.ID,
			UniKey:     item.UniKey,
			SessionId:  item.SessionId,
			FromUserId: item.SendUserId,
			Timestamp:  item.Timestamp,
			Body:       item.Body,
		})
	}
	return &logic.ListNewMessageResponse{
		Data: data,
	}, nil
}

func (l LogicGrpcServer) ListNewApplications(ctx context.Context, request *logic.ListNewApplicationsRequest) (*logic.ListNewApplicationsResponse, error) {
	applications, err := l.srv.FriendApplicationService.ListNewApplications(request.UserId)
	if err != nil {
		return &logic.ListNewApplicationsResponse{}, nil
	}
	data := make([]*logic.ApplicationItem, 0, len(applications))
	for _, app := range applications {
		data = append(data, &logic.ApplicationItem{
			Id:          app.ID,
			UserId:      app.UserId,
			Markup:      app.Markup,
			Type:        int32(app.Type),
			Status:      int32(app.Status),
			UpdatedTime: app.UpdatedTime,
		})
	}
	return &logic.ListNewApplicationsResponse{
		Data: data,
	}, nil
}

func (l LogicGrpcServer) Connect(ctx context.Context, request *logic.ConnectionRequest) (*logic.ConnectionResponse, error) {
	err := l.srv.UserService.Connect(ctx, request.UserId)
	return &logic.ConnectionResponse{}, err
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
