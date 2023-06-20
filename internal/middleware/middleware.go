package middleware

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/ykds/zura/pkg/log"
)

func InterceptorLogger(l log.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		switch level {
		case logging.LevelDebug:
			l.Debug(msg, fields)
		case logging.LevelInfo:
			l.Info(msg, fields)
		case logging.LevelWarn:
			l.Warn(msg, fields)
		case logging.LevelError:
			l.Error(msg, fields)
		}
	})
}
