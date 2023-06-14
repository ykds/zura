package friend_application

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/middleware"
)

func RegisterFriendApplicationRouter(r gin.IRouter) {
	group := r.Group("/friend-application", middleware.Auth())
	v1 := group.Group("/v1")
	{
		v1.POST("/apply", Apply)
		v1.GET("/list", ListApplications)
		v1.PUT("/id/:id", UpdateApplicationStatus)
		v1.DELETE("/id/:id", DeleteApplication)
	}
}
