package message

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/middleware"
)

func RegisterMessageRouter(r gin.IRouter) {
	group := r.Group("/message", middleware.Auth())
	{
		group.GET("/history", ListHistoryMessage)
		group.GET("/new", ListNewMessage)
		group.POST("", PushMessage)
	}
}
