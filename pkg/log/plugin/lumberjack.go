package plugin

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
)

type LumberjackConfig struct {
	Filename   string `json:"filename" yaml:"filename"`
	MaxAge     int    `json:"max_age" yaml:"max_age"`
	MaxSize    int    `json:"max_size" yaml:"max_size"`
	MaxBackups int    `json:"max_backups" yaml:"max_backups"`
	Compress   bool   `json:"compress" yaml:"compress"`
}

func NewLumberjackLogger(c LumberjackConfig) io.Writer {
	return &lumberjack.Logger{
		Filename:   c.Filename,
		MaxSize:    c.MaxSize,
		MaxAge:     c.MaxAge,
		Compress:   c.Compress,
		MaxBackups: c.MaxBackups,
		LocalTime:  true,
	}
}
