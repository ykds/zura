package session

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/middleware"
)

func RegisterSessionRouter(r gin.IRouter) {
	group := r.Group("session", middleware.Auth())
	{
		group.POST("/open", CreateSession)
		group.GET("/list", ListSession)
		group.DELETE("/id/:session_id", DeleteSession)
		group.PUT("/id/:session_id", UpdateSession)
	}
}
