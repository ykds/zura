package group

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/middleware"
)

func RegisterGroupRouter(r gin.IRouter) {
	group := r.Group("group", middleware.Auth())
	v1 := group.Group("/v1")
	{
		v1.POST("", CreateGroup)
		v1.GET("/list", ListGroup)
		v1.GET("/search", SearchGroup)
		v1.PUT("/info", UpdateGroup)
		v1.DELETE("", DismissGroup)
		v1.POST("/member", AddGroupMember)
		v1.DELETE("/member", RemoveGroupMember)
		v1.PUT("/member/info", UpdateMemberInfo)
		v1.GET("member/list", ListGroupMembers)
		v1.PUT("member/role", ChangeMemberRole)
	}
}
