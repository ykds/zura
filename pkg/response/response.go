package response

import (
	"fmt"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/log"
	"go.opentelemetry.io/otel/trace"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func HttpResponse(ctx *gin.Context, err error, data interface{}) {
	if err != nil {
		sctx := trace.SpanContextFromContext(ctx)
		tracId := sctx.TraceID()
		spanId := sctx.SpanID()
		log.Errorw(fmt.Sprintf("%+v", err), "trace_id", tracId, "span_id", spanId)
		var e errors.Error
		if errors.As(err, &e) {
			ctx.JSON(http.StatusOK, Resp{Code: e.Code.Status, Message: e.Code.Message})
		} else {
			c := errors.GetCode(errors.InternalErrorStatus)
			ctx.JSON(http.StatusOK, Resp{Code: c.Status, Message: c.Message})
		}
	} else {
		ok := errors.GetCode(errors.OKStatus)
		ctx.JSON(http.StatusOK, Resp{Code: ok.Status, Message: ok.Message, Data: data})
	}
}

func GetResponse(err error, data interface{}) Resp {
	if err != nil {
		log.Errorf("%+v", err)
		var e errors.Error
		if errors.As(err, &e) {
			return Resp{Code: e.Code.Status, Message: e.Code.Message}
		} else {
			c := errors.GetCode(errors.InternalErrorStatus)
			return Resp{Code: c.Status, Message: c.Message}
		}
	} else {
		ok := errors.GetCode(errors.OKStatus)
		return Resp{Code: ok.Status, Message: ok.Message, Data: data}
	}
}
