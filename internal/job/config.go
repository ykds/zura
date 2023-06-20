package job

type Config struct {
	Grpc        GrpcConfig        `json:"grpc" yaml:"grpc"`
	CometServer CometServerConfig `json:"comet_server" yaml:"comet_server"`
	Kafka       KafkaConfig       `json:"kafka" yaml:"kafka"`
}

type GrpcConfig struct {
	Port string `json:"port" yaml:"port"`
}

type CometServerConfig struct {
	Host string `json:"host" yaml:"host"`
	Port string `json:"port" yaml:"port"`
}

type KafkaConfig struct {
	Brokers       []string          `json:"brokers" yaml:"brokers"`
	GroupTopicMap map[string]string `json:"group_topic_map" yaml:"group_topic_map"`
}
