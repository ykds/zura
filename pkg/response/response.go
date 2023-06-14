package response

import (
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func HttpResponse(ctx *gin.Context, err error, data interface{}) {
	if err != nil {
		log.Errorf("%+v", err)
		var e errors.Error
		if errors.As(err, &e) {
			ctx.JSON(http.StatusOK, resp{Code: e.Code.Status, Message: e.Code.Message})
		} else {
			c := errors.GetCode(errors.InternalErrorStatus)
			ctx.JSON(http.StatusOK, resp{Code: c.Status, Message: c.Message})
		}
	} else {
		ok := errors.GetCode(errors.OKStatus)
		ctx.JSON(http.StatusOK, resp{Code: ok.Status, Message: ok.Message, Data: data})
	}
}
