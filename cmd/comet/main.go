package main

import (
	"flag"
	"github.com/ykds/zura/internal/comet"
	cfg "github.com/ykds/zura/pkg/config"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/pkg/log/zap"
	"github.com/ykds/zura/pkg/snowflake"
	"os"
	"os/signal"
	"syscall"
)

var configPath = flag.String("conf", "./config.yaml", "config file path")

func main() {
	flag.Parse()
	cfg.InitConfig(*configPath, comet.GetConfig())

	snowflake.InitSnowflake(2)

	l := zap.NewLogger(&comet.GetConfig().Log,
		zap.WithDebug(comet.GetConfig().Debug),
		zap.WithLumberjack())
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
