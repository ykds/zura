package log

type Config struct {
	Level      string `json:"level" yaml:"level"`
	Filename   string `json:"filename" yaml:"filename"`
	MaxAge     int    `json:"max_age" yaml:"max_age"`
	MaxSize    int    `json:"max_size" yaml:"max_size"`
	MaxBackups int    `json:"max_backups" yaml:"max_backups"`
	Compress   bool   `json:"compress" yaml:"compress"`
}

func DefaultConfig() Config {
	return Config{
		Level:      DebugLevel,
		Filename:   "./logs/zura.log",
		MaxAge:     7,
		MaxSize:    10,
		MaxBackups: 5,
		Compress:   false,
	}
}
