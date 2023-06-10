package friends

import (
	"zura/internal/logic/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterFriendsRouter(r gin.IRouter) {
	group := r.Group("friends", middleware.Auth())
	{
		group.GET("/list", ListFriends)
		group.DELETE("", DeleteFriends)
	}
}
