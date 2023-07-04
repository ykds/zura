package zap

import (
	"github.com/ykds/zura/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

type logger struct {
	*zap.SugaredLogger
	c       log.Config
	debug   bool
	outs    []io.Writer
	service string
}

func (l *logger) Write(p []byte) (n int, err error) {
	for _, out := range l.outs {
		_, _ = out.Write(p)
	}
	return 0, nil
}

type Option func(*logger)

func WithDebug(debug bool) Option {
	return func(l *logger) {
		l.debug = debug
	}
}

func WithOutput(writer ...io.Writer) Option {
	return func(l *logger) {
		l.outs = append(l.outs, writer...)
	}
}

func WithService(service string) Option {
	return func(l *logger) {
		l.service = service
	}
}

func NewLogger(cfg log.Config, opts ...Option) log.Logger {
	lg := &logger{
		c:    cfg,
		outs: make([]io.Writer, 0),
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
		lg.outs = []io.Writer{os.Stdout}
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	enc := zapcore.NewJSONEncoder(encoderConfig)

	wsList := make([]zapcore.WriteSyncer, 0, len(lg.outs))
	for _, out := range lg.outs {
		wsList = append(wsList, zapcore.AddSync(out))
	}

	core := zapcore.NewCore(enc, zapcore.NewMultiWriteSyncer(wsList...), level)
	l := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(l)
	lg.SugaredLogger = l.Sugar()
	lg.SugaredLogger.With("service", lg.service)
	return lg
}
