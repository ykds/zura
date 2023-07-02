package main

import (
	"context"
	"flag"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/internal/consumer/config"
	"github.com/ykds/zura/internal/consumer/logging"
	cfg "github.com/ykds/zura/pkg/config"
	"github.com/ykds/zura/pkg/kafka"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/pkg/log/plugin"
	"github.com/ykds/zura/pkg/log/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configPath = flag.String("conf", "./config.yaml", "config file path")

func main() {
	flag.Parse()
	cfg.InitConfig(*configPath, config.GetConfig())

	l := zap.NewLogger(
		config.GetConfig().Log,
		zap.WithService("logic"),
		zap.WithDebug(config.GetConfig().Debug),
		zap.WithOutput(plugin.NewLumberjackLogger(config.DefaultConfig().Log.Lumberjack), plugin.NewLogstash(config.GetConfig().Log.Logstash)))
	log.SetGlobalLogger(l)

	ctx, cancel := context.WithCancel(context.Background())

	kafkaManager := kafka.NewKafka(config.GetConfig().Kafka)

	c := kafkaManager.NewConsumer("logging-consumer-1", []string{common.LoggingTopic}, logging.NewLoggingConsumer(), kafka.WithCommitInterval(5*time.Second))
	go c.Run(ctx)

	log.Info("server started.")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sig:
		cancel()
	}
	log.Info("exit.")
}
