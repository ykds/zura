package main

import (
	"flag"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/internal/logic/config"
	"github.com/ykds/zura/internal/logic/entity"
	"github.com/ykds/zura/internal/logic/server"
	"github.com/ykds/zura/internal/logic/services"
	"github.com/ykds/zura/pkg/cache"
	cfg "github.com/ykds/zura/pkg/config"
	"github.com/ykds/zura/pkg/db"
	"github.com/ykds/zura/pkg/kafka"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/pkg/log/plugin"
	"github.com/ykds/zura/pkg/log/zap"
	"github.com/ykds/zura/pkg/snowflake"
	"github.com/ykds/zura/pkg/trace"
	"os"
	"os/signal"
	"syscall"
)

var configPath = flag.String("conf", "./config.yaml", "config file path")

func main() {
	flag.Parse()
	cfg.InitConfig(*configPath, config.GetConfig())

	snowflake.InitSnowflake(1)
	trace.InitTrace(config.GetConfig().Trace)

	kafkaManager := kafka.NewKafka(config.GetConfig().Kafka)
	defer kafkaManager.Close()

	producer := kafkaManager.NewProducer(common.LoggingTopic)

	l := zap.NewLogger(
		config.GetConfig().Log,
		zap.WithService("logic"),
		zap.WithDebug(config.GetConfig().Server.Debug),
		zap.WithOutput(os.Stdout, plugin.NewKafkaWriter(producer)))
	log.SetGlobalLogger(l)

	database := db.New(&config.GetConfig().Database, db.WithDebug(config.GetConfig().Server.Debug))
	redis := cache.NewRedis(config.GetConfig().Cache)
	cache.SetGlobalCache(redis)
	entity.NewEntity(database, redis)

	services.NewServices(redis, entity.GetEntity(), kafkaManager)

	httpServer := server.NewHttpServer(config.GetConfig().HttpServer,
		server.WithLogger(l),
		server.WithDebug(config.GetConfig().Server.Debug))
	httpServer.Run()

	logicGrpcSrv := server.NewGrpcServer(config.GetConfig().GrpcServer, *services.GetServices())

	log.Info("server started.")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sig:
		logicGrpcSrv.GracefulStop()
		_ = httpServer.Shutdown()
	}
	log.Info("exit.")
}
