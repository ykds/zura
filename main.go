package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"zura/internal/config"
	"zura/internal/entity"
	"zura/internal/server"
	"zura/internal/services"
	"zura/pkg/cache"
	"zura/pkg/db"
	"zura/pkg/log"
	"zura/pkg/log/zap"
)

var configPath = flag.String("conf", "./config.yaml", "config file path")

func main() {
	flag.Parse()
	config.InitConfig(*configPath)

	l := zap.NewLogger(&config.GetConfig().Log, 
		zap.WithDebug(config.GetConfig().Server.Debug), 
		zap.WithLumberjack())
	log.SetGlobalLogger(l)

	ctx, cancel := context.WithCancel(context.Background())

	database := db.New(&config.GetConfig().Database, db.WithDebug(config.GetConfig().Server.Debug))
	cache := cache.NewRedis(ctx, &config.GetConfig().Cache)
	entity.NewEntity(database, cache)
	services.NewServices(cache, entity.GetEntity())

	httpServer := server.NewHttpServer(ctx, 
		server.WithConfig(config.GetConfig().HttpServer), 
		server.WithLogger(l), 
		server.WithDebug(config.GetConfig().Server.Debug))
	httpServer.Run()
	log.Info("server started.")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sig:
		cancel()
	}
	log.Info("exit.")
}
