package config

import "zira/pkg/db"

var cfg *Config

type Config struct {
	Database db.Config `json:"database" yaml:"database"`
}

func GetConfig() *Config {
	return cfg
}

func InitConfig() error {
	cfg = &Config{}
	return nil
}

