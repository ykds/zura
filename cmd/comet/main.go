package main

import (
	"flag"
	"github.com/ykds/zura/internal/comet"
	"github.com/ykds/zura/internal/common"
	cfg "github.com/ykds/zura/pkg/config"
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
	cfg.InitConfig(*configPath, comet.GetConfig())

	snowflake.InitSnowflake(2)
	trace.InitTrace(comet.GetConfig().Trace)

	kafkaManager := kafka.NewKafka(comet.GetConfig().Kafka)
	producer := kafkaManager.NewProducer(common.LoggingTopic)

	l := zap.NewLogger(
		comet.GetConfig().Log,
		zap.WithService("comet"),
		zap.WithDebug(comet.GetConfig().Debug),
		zap.WithOutput(plugin.NewLumberjackLogger(comet.GetConfig().Log.Lumberjack), plugin.NewKafkaWriter(producer)))
	log.SetGlobalLogger(l)

	server := comet.NewServer(comet.GetConfig())
	grpcServer := comet.NewGrpcServer(server)

	log.Info("server started.")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sig:
		grpcServer.GracefulStop()
	}
	log.Info("exit.")
}
