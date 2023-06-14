package middleware

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"zura/pkg/log"
)

func InterceptorLogger(l log.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		switch level {
		case logging.LevelDebug:
			l.Infof(msg, fields...)
		case logging.LevelInfo:
			l.Infof(msg, fields...)
		case logging.LevelWarn:
			l.Warnf(msg, fields...)
		case logging.LevelError:
			l.Errorf(msg, fields...)
		}
	})
}
