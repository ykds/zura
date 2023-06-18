package common

import (
	"github.com/gin-gonic/gin"
	"net/url"
)

var (
	RegisterVerifyCodeKey = "register_verify_code:%s:%s"

	StaticDir  = "./static/"
	StaticPath = "/static/"

	UserIdKey = "userId"
)

func ParseAvatarUrl(c *gin.Context, avatar string) string {
	if avatar != "" {
		avatarUrl := url.URL{
			Scheme: "http",
			Host:   c.Request.Host,
			Path:   StaticPath + avatar,
		}
		return avatarUrl.String()
	}
	return ""
}
