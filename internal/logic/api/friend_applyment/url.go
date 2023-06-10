package friend_applyment

import (
	"zura/internal/logic/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterFriendApplymentRouter(r gin.IRouter) {
	group := r.Group("/friend_applyment", middleware.Auth())
	{
		group.POST("/apply", Apply)
		group.GET("/list", ListApplyments)
		group.PUT("/id/:id", UpdateApplymentStatus)
		group.DELETE("/id/:id", DeleteApplyment)
	}
}
