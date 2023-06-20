package comet

import "github.com/ykds/zura/pkg/log"

var cfg = DefaultConfig()

func GetConfig() *Config {
	return cfg
}

type Config struct {
	Debug    bool       `json:"debug" yaml:"debug"`
	HttpPort string     `json:"http_port" yaml:"http_port"`
	GrpcPort string     `json:"grpc_port" yaml:"grpc_port"`
	Logic    Logic      `json:"logic" yaml:"logic"`
	Log      log.Config `json:"log" yaml:"log"`
	Session  Session    `json:"session" yaml:"session"`
}

func DefaultConfig() *Config {
	return &Config{
		Debug:    true,
		HttpPort: "9000",
		GrpcPort: "9001",
		Logic: Logic{
			Host: "localhost",
			Port: "8001",
		},
		Log: log.DefaultConfig(),
		Session: Session{
			HeartbeatInterval: 30,
		},
	}
}

type Logic struct {
	Host string `json:"host" yaml:"host"`
	Port string `json:"port" yaml:"port"`
}

type Session struct {
	HeartbeatInterval int8 `json:"heartbeat_interval" yaml:"heartbeat_interval"`
}
