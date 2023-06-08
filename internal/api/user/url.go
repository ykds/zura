package user

import (
	"zura/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r gin.IRouter) {
	noAuth := r.Group("/users")
	{
		noAuth.POST("/register", Register)
		noAuth.POST("/login", Login)
	}
	auth := r.Group("/users", middleware.Auth())
	{
		auth.GET("/info", GetUserInfo)
	}
}
