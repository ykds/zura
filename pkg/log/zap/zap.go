package zap

import (
	"github.com/ykds/zura/pkg/log"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type logger struct {
	*zap.SugaredLogger
	c     *log.Config
	debug bool
	out   io.Writer
}

func (l *logger) Write(p []byte) (n int, err error) {
	return l.out.Write(p)
}

type Option func(*logger)

func WithDebug(debug bool) Option {
	return func(l *logger) {
		l.debug = debug
	}
}

func WithLumberjack() Option {
	return func(l *logger) {
		l.out = &lumberjack.Logger{
			Filename:   l.c.Filename,
			MaxSize:    l.c.MaxSize,
			MaxAge:     l.c.MaxAge,
			Compress:   l.c.Compress,
			MaxBackups: l.c.MaxBackups,
			LocalTime:  true,
		}
	}
}

func NewLogger(cfg *log.Config, opts ...Option) log.Logger {
	lg := &logger{
		c: cfg,
	}
	for _, opt := range opts {
		opt(lg)
	}

	level := zap.InfoLevel
	switch cfg.Level {
	case log.DebugLevel:
		level = zap.DebugLevel
	case log.WarnLevel:
		level = zap.WarnLevel
	case log.ErrorLevel:
		level = zap.ErrorLevel
	}

	if lg.debug {
		level = zap.DebugLevel
		lg.out = os.Stdout
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	enc := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(enc, zapcore.AddSync(lg.out), level)
	l := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(l)
	lg.SugaredLogger = l.Sugar()
	return lg
}
