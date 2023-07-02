package kafka

type Config struct {
	Brokers []string `json:"brokers" yaml:"brokers"`
}

func DefaultConfig() Config {
	return Config{
		Brokers: []string{"localhost:9092"},
	}
}
