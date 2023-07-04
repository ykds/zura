package middleware

import (
	"context"
	"fmt"
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
			l.Debugw(fmt.Sprintf("%s, field: %v", msg, fields), []any{"trace_id", tracId.String(), "span_id", spanId.String()}...)
		case logging.LevelInfo:
			l.Infow(fmt.Sprintf("%s, field: %v", msg, fields), []any{"trace_id", tracId.String(), "span_id", spanId.String()}...)
		case logging.LevelWarn:
			l.Warnw(fmt.Sprintf("%s, field: %v", msg, fields), []any{"trace_id", tracId.String(), "span_id", spanId.String()}...)
		case logging.LevelError:
			l.Errorw(fmt.Sprintf("%s, field: %v", msg, fields), []any{"trace_id", tracId.String(), "span_id", spanId.String()}...)
		}
	})
}
