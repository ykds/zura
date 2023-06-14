package session

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/middleware"
)

func RegisterSessionRouter(r gin.IRouter) {
	group := r.Group("session", middleware.Auth())
	{
		group.POST("/open", OpenSession)
		group.GET("/list", ListSession)
		group.DELETE("/id/:session_id", DeleteSession)
	}
}
