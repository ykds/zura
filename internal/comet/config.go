package comet

import (
	"github.com/ykds/zura/pkg/kafka"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/pkg/trace"
)

var cfg = DefaultConfig()

func GetConfig() *Config {
	return cfg
}

type Config struct {
	Debug      bool             `json:"debug" yaml:"debug"`
	HttpServer HttpServerConfig `json:"http_server" yaml:"http_server"`
	GrpcServer GrpcServerConfig `json:"grpc_server" yaml:"grpc_server"`
	Logic      Logic            `json:"logic" yaml:"logic"`
	Log        log.Config       `json:"log" yaml:"log"`
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

func DefaultConfig() *Config {
	return &Config{
		Debug: true,
		HttpServer: HttpServerConfig{
			Port: "9080",
		},
		GrpcServer: GrpcServerConfig{
			Port: "9001",
		},
		Logic: Logic{
			Host: "localhost",
			Port: "8001",
		},
		Log: log.DefaultConfig(),
		Session: Session{
			HeartbeatInterval: 30,
		},
		Trace: trace.Config{
			ServiceName: "comet",
		},
		Kafka: kafka.DefaultConfig(),
	}
}

type Logic struct {
	Host string `json:"host" yaml:"host"`
	Port string `json:"port" yaml:"port"`
}

type Session struct {
	HeartbeatInterval int8 `json:"heartbeat_interval" yaml:"heartbeat_interval"`
}
