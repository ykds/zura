package log

type Config struct {
	Level string `json:"level" yaml:"level"`
}

func DefaultConfig() Config {
	return Config{
		Level: DebugLevel,
	}
}
