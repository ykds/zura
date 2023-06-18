package group

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/middleware"
)

func RegisterGroupRouter(r gin.IRouter) {
	group := r.Group("group", middleware.Auth())
	{
		group.POST("", CreateGroup)
		group.GET("/list", ListGroup)
		group.PUT("/info", UpdateGroup)
		group.DELETE("", DismissGroup)
		group.POST("/member", AddGroupMember)
		group.DELETE("/member", RemoveGroupMember)
		group.PUT("/member/info", UpdateMemberInfo)
		group.GET("member/list", ListGroupMembers)
		group.PUT("member/role", ChangeMemberRole)
	}
}
