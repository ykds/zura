package main

import (
	"context"
	"flag"
	"github.com/ykds/zura/internal/logic/config"
	"github.com/ykds/zura/internal/logic/entity"
	"github.com/ykds/zura/internal/logic/server"
	"github.com/ykds/zura/internal/logic/services"
	"github.com/ykds/zura/pkg/cache"
	cfg "github.com/ykds/zura/pkg/config"
	"github.com/ykds/zura/pkg/db"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/pkg/log/zap"
	"github.com/ykds/zura/pkg/snowflake"
	"github.com/ykds/zura/proto/comet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configPath = flag.String("conf", "./config.yaml", "config file path")

func main() {
	flag.Parse()
	cfg.InitConfig(*configPath, config.GetConfig())

	snowflake.InitSnowflake()

	l := zap.NewLogger(&config.GetConfig().Log,
		zap.WithDebug(config.GetConfig().Server.Debug),
		zap.WithLumberjack())
	log.SetGlobalLogger(l)

	ctx, cancel := context.WithCancel(context.Background())

	database := db.New(&config.GetConfig().Database, db.WithDebug(config.GetConfig().Server.Debug))
	caches := cache.NewRedis(ctx, &config.GetConfig().Cache)
	entity.NewEntity(database, caches)

	ctx2, cancel2 := context.WithTimeout(ctx, 2*time.Second)
	cometConn, err := grpc.DialContext(ctx2, "127.0.0.1:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panicf("new comet grpc client failed, err: %+v", err)
	}
	cancel2()
	cometClient := comet.NewCometClient(cometConn)

	services.NewServices(caches, entity.GetEntity(), cometClient)

	httpServer := server.NewHttpServer(ctx,
		server.WithConfig(config.GetConfig().HttpServer),
		server.WithLogger(l),
		server.WithDebug(config.GetConfig().Server.Debug))
	httpServer.Run()

	logicGrpcSrv := server.NewGrpcServer(*services.GetServices())

	log.Info("server started.")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sig:
		logicGrpcSrv.GracefulStop()
		cancel()
	}
	log.Info("exit.")
}
