package config

import (
	"github.com/ykds/zura/pkg/cache"
	"github.com/ykds/zura/pkg/db"
	"github.com/ykds/zura/pkg/kafka"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/pkg/trace"
)

var cfg = DefaultConfig()

type ServerConfig struct {
	Debug bool `json:"debug" yaml:"debug"`
}

type Config struct {
	Server     ServerConfig     `json:"server" yaml:"server"`
	Database   db.Config        `json:"database" yaml:"database"`
	Cache      cache.Config     `json:"cache" yaml:"cache"`
	Log        log.Config       `json:"log" yaml:"log"`
	HttpServer HttpServerConfig `json:"http_server" yaml:"http_server"`
	GrpcServer GrpcServerConfig `json:"grpc_server" yaml:"grpc_server"`
	Session    Session          `json:"session" yaml:"session"`
	Trace      trace.Config     `json:"trace" yaml:"trace"`
	Kafka      kafka.Config     `json:"kafka" yaml:"kafka"`
}

type HttpServerConfig struct {
	Port string `json:"port"`
}

type GrpcServerConfig struct {
	Port string `json:"port" yaml:"port"`
}

type Session struct {
	HeartbeatInterval int `json:"heartbeat_interval" yaml:"heartbeat_interval"`
}

func DefaultConfig() *Config {
	return &Config{
		Database: db.DefaultConfig(),
		Cache:    cache.DefaultConfig(),
		Log:      log.DefaultConfig(),
		HttpServer: HttpServerConfig{
			Port: "8080",
		},
		GrpcServer: GrpcServerConfig{
			Port: "8001",
		},
		Session: Session{
			HeartbeatInterval: 60,
		},
		Trace: trace.Config{
			ServiceName: "logic",
		},
		Kafka: kafka.DefaultConfig(),
	}
}

func GetConfig() *Config {
	return cfg
}
