package config

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"zura/internal/logic/server"
	"zura/pkg/cache"
	"zura/pkg/db"
	"zura/pkg/log"

	"gopkg.in/yaml.v2"
)

var cfg *Config

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

func DefautConfig() *Config {
	return &Config{
		Database:   db.DefaultConfig(),
		Cache:      cache.DefaultConfig(),
		Log:        log.DefaultConfig(),
		HttpServer: server.DefaultConfig(),
	}
}

func GetConfig() *Config {
	if cfg == nil {
		panic("未初始化配置")
	}
	return cfg
}

func InitConfig(path string) {
	cfg = DefautConfig()
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if strings.HasSuffix(path, ".json") {
		err = json.Unmarshal(b, cfg)
		if err != nil {
			panic(err)
		}
	}
	if strings.HasSuffix(path, ".yaml") {
		err = yaml.Unmarshal(b, cfg)
		if err != nil {
			panic(err)
		}
	}
}
