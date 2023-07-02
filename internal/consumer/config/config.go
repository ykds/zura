package config

import (
	"github.com/ykds/zura/pkg/kafka"
	"github.com/ykds/zura/pkg/log"
)

var cfg = DefaultConfig()

type Config struct {
	Debug bool         `json:"debug" yaml:"debug"`
	Kafka kafka.Config `json:"kafka" yaml:"kafka"`
	Log   log.Config   `json:"log" yaml:"log"`
}

func DefaultConfig() *Config {
	return &Config{
		Kafka: kafka.DefaultConfig(),
	}
}

func GetConfig() *Config {
	return cfg
}
