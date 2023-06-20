package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/middleware"
)

func RegisterUserRouter(r gin.IRouter) {
	noAuth := r.Group("/users")
	v1noAuth := noAuth.Group("/v1")
	{
		v1noAuth.POST("/register", Register)
		v1noAuth.POST("/token", Login)
	}
	auth := r.Group("/users", middleware.Auth())
	v1auth := auth.Group("/v1")
	{
		v1auth.GET("/info", GetUserInfo)
		v1auth.PUT("/info", UpdateInfo)
		v1auth.GET("/search", SearchUser)
	}
}
