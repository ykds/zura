package session

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/middleware"
)

func RegisterSessionRouter(r gin.IRouter) {
	group := r.Group("session", middleware.Auth())
	v1 := group.Group("/v1")
	{
		v1.POST("/open", CreateSession)
		v1.GET("/list", ListSession)
		v1.DELETE("/id/:session_id", DeleteSession)
		v1.PUT("/id/:session_id", UpdateSession)
	}
}
