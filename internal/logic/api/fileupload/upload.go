package fileupload

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/pkg/response"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"time"
)

func Upload(c *gin.Context) {
	var (
		err  error
		file *multipart.FileHeader
		resp struct {
			Url      string `json:"url"`
			Filename string `json:"filename"`
		}
	)
	defer func() {
		response.HttpResponse(c, err, resp)
	}()
	file, err = c.FormFile("file")
	if err != nil {
		return
	}
	ext := filepath.Ext(file.Filename)
	file.Filename = fmt.Sprintf("%d%s", time.Now().UnixMilli(), ext)
	err = c.SaveUploadedFile(file, common.StaticDir+file.Filename)
	if err != nil {
		return
	}
	resp.Url = (&url.URL{Scheme: "http", Host: c.Request.Host, Path: common.StaticPath + file.Filename}).String()
	resp.Filename = file.Filename
}
