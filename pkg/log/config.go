package log

import (
	"github.com/ykds/zura/pkg/log/plugin"
)

type Config struct {
	Level      string                  `json:"level" yaml:"level"`
	Lumberjack plugin.LumberjackConfig `json:"lumberjack" yaml:"lumberjack"`
	Logstash   plugin.LogstashConfig   `json:"logstash" yaml:"logstash"`
}

func DefaultConfig() Config {
	return Config{
		Level: DebugLevel,
		Lumberjack: plugin.LumberjackConfig{
			Filename:   "./logs/zura.log",
			MaxAge:     7,
			MaxSize:    10,
			MaxBackups: 5,
			Compress:   false,
		},
		Logstash: plugin.LogstashConfig{
			Host: "localhost",
			Port: "4560",
		},
	}
}
