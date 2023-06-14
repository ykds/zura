package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/middleware"
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
		auth.POST("/upload", UploadUserAvatar)
		auth.GET("/search", SearchUser)
	}
}
