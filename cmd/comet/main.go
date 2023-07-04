package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ykds/zura/internal/comet"
	"github.com/ykds/zura/internal/common"
	cfg "github.com/ykds/zura/pkg/config"
	"github.com/ykds/zura/pkg/discovery"
	"github.com/ykds/zura/pkg/kafka"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/pkg/log/plugin"
	"github.com/ykds/zura/pkg/log/zap"
	"github.com/ykds/zura/pkg/net"
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

	snowflake.InitSnowflake(int64(comet.GetConfig().Server.ID))
	trace.InitTrace(comet.GetConfig().Trace)

	appCtx, cancel := context.WithCancel(context.Background())

	kafkaManager := kafka.NewKafka(comet.GetConfig().Kafka)
	defer kafkaManager.Close()

	producer := kafkaManager.NewProducer(common.LoggingTopic)

	l := zap.NewLogger(
		comet.GetConfig().Log,
		zap.WithService("comet"),
		zap.WithDebug(comet.GetConfig().Server.Debug),
		zap.WithOutput(os.Stdout, plugin.NewKafkaWriter(producer)))
	log.SetGlobalLogger(l)

	server := comet.NewServer(comet.GetConfig())
	grpcServer := comet.NewGrpcServer(server)

	etcd := discovery.NewEtcd(comet.GetConfig().Etcd, common.CometDiscoveryEndpoint)
	ip, err := net.GetOutBoundIP()
	if err != nil {
		return
	}
	go func() {
		err2 := etcd.Register(appCtx, fmt.Sprintf("/%d", comet.GetConfig().Server.ID), ip+":"+comet.GetConfig().GrpcServer.Port, nil)
		if err2 != nil {
			log.Panicf("register comet failed: %v", err2)
		}
	}()
	log.Info("server started.")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sig:
		grpcServer.GracefulStop()
		cancel()
	}
	log.Info("exit.")
}
