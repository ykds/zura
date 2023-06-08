package middleware

import (
	"zura/internal/codec"
	"zura/pkg/errors"
	"zura/pkg/response"
	"zura/pkg/token"

	"github.com/gin-gonic/gin"
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
		ctx.Set("userId", userId)
		ctx.Next()
	}
}
