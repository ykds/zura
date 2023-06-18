package fileupload

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/middleware"
)

func RegisterUploadRouter(r gin.IRouter) {
	group := r.Group("/upload", middleware.Auth())
	{
		group.POST("", Upload)
	}
}
