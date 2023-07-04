package common

import (
	"github.com/gin-gonic/gin"
	"net/url"
)

const (
	RegisterVerifyCodeKey = "register_verify_code:%s:%s"

	StaticDir  = "./static/"
	StaticPath = "/static/"

	UserIdKey = "userId"

	MessageCacheKey      = "message_%d"
	GroupMessageCacheKey = "group_message_%d"
	UserOnlineCacheKey   = "CACHE_ONLINE_USER:%d"

	UpdateUserPhoneCacheKey = "UPDATE_USER_PHONE_%d"
	UpdateUserEmailCacheKey = "UPDATE_USER_EMAIL_%d"

	CometDiscoveryEndpoint = "/zura/comet"
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
