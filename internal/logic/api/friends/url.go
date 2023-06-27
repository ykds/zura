package friends

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/middleware"
)

func RegisterFriendsRouter(r gin.IRouter) {
	group := r.Group("friends", middleware.Auth())
	v1 := group.Group("/v1")
	{
		v1.GET("/list", ListFriends)
		v1.DELETE("/id/:id", DeleteFriends)
	}
}
