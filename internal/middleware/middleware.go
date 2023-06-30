package middleware

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/ykds/zura/pkg/log"
	"go.opentelemetry.io/otel/trace"
)

func InterceptorLogger(l log.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		sctx := trace.SpanContextFromContext(ctx)
		tracId := sctx.TraceID()
		spanId := sctx.SpanID()
		switch level {
		case logging.LevelDebug:
			l.Debugw(msg, []any{"field", fields, "trace_id", tracId, "span_id", spanId}...)
		case logging.LevelInfo:
			l.Infow(msg, []any{"field", fields, "trace_id", tracId, "span_id", spanId}...)
		case logging.LevelWarn:
			l.Warnw(msg, []any{"field", fields, "trace_id", tracId, "span_id", spanId}...)
		case logging.LevelError:
			l.Errorw(msg, []any{"field", fields, "trace_id", tracId, "span_id", spanId}...)
		}
	})
}
