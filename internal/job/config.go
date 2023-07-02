package job

import "github.com/ykds/zura/pkg/kafka"

type Config struct {
	Grpc        GrpcConfig        `json:"grpc" yaml:"grpc"`
	CometServer CometServerConfig `json:"comet_server" yaml:"comet_server"`
	Kafka       kafka.Config      `json:"kafka" yaml:"kafka"`
}

type GrpcConfig struct {
	Port string `json:"port" yaml:"port"`
}

type CometServerConfig struct {
	Host string `json:"host" yaml:"host"`
	Port string `json:"port" yaml:"port"`
}
