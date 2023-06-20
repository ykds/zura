package message

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/middleware"
)

func RegisterMessageRouter(r gin.IRouter) {
	group := r.Group("/message", middleware.Auth())
	v1 := group.Group("/v1")
	{
		v1.GET("/history", ListHistoryMessage)
		//v1.GET("/new", ListNewMessage)
		v1.POST("/push", PushMessage)
	}
}
