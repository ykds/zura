package friends

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/middleware"
)

func RegisterFriendsRouter(r gin.IRouter) {
	group := r.Group("friends", middleware.Auth())
	{
		group.GET("/list", ListFriends)
		group.DELETE("", DeleteFriends)
	}
}
