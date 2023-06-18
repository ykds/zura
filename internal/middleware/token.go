package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/internal/logic/codec"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/response"
	"github.com/ykds/zura/pkg/token"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := ctx.GetHeader("token")
		if t == "" {
			response.HttpResponse(ctx, errors.New(codec.NeedAuthStatus), nil)
			ctx.Abort()
			return
		}
		userId, err := token.VerifyToken(t)
		if err != nil {
			response.HttpResponse(ctx, errors.WithMessage(errors.New(codec.ParseTokenFailedStatus), err.Error()), nil)
			ctx.Abort()
			return
		}
		ctx.Set(common.UserIdKey, userId)
		ctx.Next()
	}
}
