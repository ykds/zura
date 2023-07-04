package main

import (
	"context"
	"flag"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/internal/consumer/config"
	"github.com/ykds/zura/internal/consumer/logging"
	"github.com/ykds/zura/internal/consumer/message"
	"github.com/ykds/zura/internal/consumer/notify"
	cfg "github.com/ykds/zura/pkg/config"
	"github.com/ykds/zura/pkg/discovery"
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
	logstash := plugin.NewLogstash(config.GetConfig().Logstash)
	defer logstash.Close()
	l := zap.NewLogger(
		config.GetConfig().Log,
		zap.WithService("logic"),
		zap.WithDebug(config.GetConfig().Debug),
		zap.WithOutput(os.Stdout, logstash))
	log.SetGlobalLogger(l)

	ctx, cancel := context.WithCancel(context.Background())

	kafkaManager := kafka.NewKafka(config.GetConfig().Kafka)
	defer kafkaManager.Close()

	c := kafkaManager.NewConsumer("logging-consumer-1", []string{common.LoggingTopic}, logging.NewLoggingConsumer(), kafka.WithCommitInterval(5*time.Second))
	go c.Run(ctx)

	dis := discovery.NewEtcd(config.GetConfig().Etcd, common.CometDiscoveryEndpoint)
	defer dis.Close()

	c1 := kafkaManager.NewConsumer("message-consumer-1", []string{common.MessageTopic}, message.NewConsumer(ctx, dis), kafka.WithCommitInterval(5*time.Second))
	go c1.Run(ctx)

	c2 := kafkaManager.NewConsumer("notify-consumer-1", []string{common.NotificationTopic}, notify.NewConsumer(ctx, dis), kafka.WithCommitInterval(5*time.Second))
	go c2.Run(ctx)

	log.Info("server started.")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sig:
		cancel()
	}
	log.Info("exit.")
}
