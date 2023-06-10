package session

import (
	"zura/internal/logic/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterSessionRouter(r gin.IRouter) {
	group := r.Group("session", middleware.Auth())
	{
		group.POST("/open", OpenSession)
		group.GET("/list", ListSession)
		group.DELETE("/id/:session_id", DeleteSession)
	}
}
