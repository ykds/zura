package config

import (
	"github.com/ykds/zura/pkg/discovery"
	"github.com/ykds/zura/pkg/kafka"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/pkg/log/plugin"
)

var cfg = DefaultConfig()

type Config struct {
	Debug    bool                  `json:"debug" yaml:"debug"`
	Kafka    kafka.Config          `json:"kafka" yaml:"kafka"`
	Log      log.Config            `json:"log" yaml:"log"`
	Etcd     discovery.Config      `json:"etcd" yaml:"etcd"`
	Logstash plugin.LogstashConfig `json:"logstash" yaml:"logstash"`
}

func DefaultConfig() *Config {
	return &Config{
		Kafka: kafka.DefaultConfig(),
		Etcd: discovery.Config{
			Urls: []string{"http://localhost:2379"},
		},
	}
}

func GetConfig() *Config {
	return cfg
}
