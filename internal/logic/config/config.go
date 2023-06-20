package config

import (
	"github.com/ykds/zura/internal/logic/server"
	"github.com/ykds/zura/pkg/cache"
	"github.com/ykds/zura/pkg/db"
	"github.com/ykds/zura/pkg/log"
)

var cfg = DefaultConfig()

type ServerConfig struct {
	Debug bool `json:"debug" yaml:"debug"`
}

type Config struct {
	Server      ServerConfig            `json:"server" yaml:"server"`
	Database    db.Config               `json:"database" yaml:"database"`
	Cache       cache.Config            `json:"cache" yaml:"cache"`
	Log         log.Config              `json:"log" yaml:"log"`
	HttpServer  server.HttpServerConfig `json:"http_server" yaml:"http_server"`
	CometServer struct {
		Host string `json:"host" yaml:"host"`
		Port string `json:"port" yaml:"port"`
	} `json:"comet_server" yaml:"comet_server"`
	GrpcServer server.GrpcServerConfig `json:"grpc_server" yaml:"grpc_server"`
	Session    Session                 `json:"session" yaml:"session"`
}

type Session struct {
	HeartbeatInterval int8 `json:"heartbeat_interval" yaml:"heartbeat_interval"`
}

func DefaultConfig() *Config {
	return &Config{
		Database:   db.DefaultConfig(),
		Cache:      cache.DefaultConfig(),
		Log:        log.DefaultConfig(),
		HttpServer: server.DefaultConfig(),
		GrpcServer: server.DefaultGrpcConfig(),
		Session: Session{
			HeartbeatInterval: 60,
		},
	}
}

func GetConfig() *Config {
	return cfg
}
