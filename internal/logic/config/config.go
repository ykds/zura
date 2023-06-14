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
	Server     ServerConfig            `json:"server" yaml:"server"`
	Database   db.Config               `json:"database" yaml:"database"`
	Cache      cache.Config            `json:"cache" yaml:"cache"`
	Log        log.Config              `json:"log" yaml:"log"`
	HttpServer server.HttpServerConfig `json:"http_server" yaml:"http_server"`
}

func DefaultConfig() *Config {
	return &Config{
		Database:   db.DefaultConfig(),
		Cache:      cache.DefaultConfig(),
		Log:        log.DefaultConfig(),
		HttpServer: server.DefaultConfig(),
	}
}

func GetConfig() *Config {
	return cfg
}
