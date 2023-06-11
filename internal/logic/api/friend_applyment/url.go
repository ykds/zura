package friend_applyment

import (
	"zura/internal/logic/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterFriendApplymentRouter(r gin.IRouter) {
	group := r.Group("/friend-applyment", middleware.Auth())
	v1 := group.Group("/v1")
	{
		v1.POST("/apply", Apply)
		v1.GET("/list", ListApplyments)
		v1.PUT("/id/:id", UpdateApplymentStatus)
		v1.DELETE("/id/:id", DeleteApplyment)
	}
}
