package message

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/internal/logic/services"
	"github.com/ykds/zura/internal/logic/services/message"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/response"
)

func ListHistoryMessage(c *gin.Context) {
	var (
		err  error
		req  message.ListMessageRequest
		resp struct {
			Data []message.MessageItem `json:"data"`
		}
	)
	defer func() {
		if len(resp.Data) == 0 {
			resp.Data = []message.MessageItem{}
		}
		response.HttpResponse(c, err, resp)
	}()
	if err = c.BindQuery(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	resp.Data, err = services.GetServices().MessageService.ListHistoryMessage(c.GetInt64(common.UserIdKey), req)
}

func ListNewMessage(c *gin.Context) {
	var (
		err  error
		req  message.ListMessageRequest
		resp struct {
			Data []message.MessageItem `json:"data"`
		}
	)
	defer func() {
		if len(resp.Data) == 0 {
			resp.Data = []message.MessageItem{}
		}
		response.HttpResponse(c, err, resp)
	}()
	if err = c.BindQuery(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	resp.Data, err = services.GetServices().MessageService.ListNewMessage(c.GetInt64(common.UserIdKey), req)
}

func PushMessage(c *gin.Context) {
	var (
		err error
		req message.PushMessageRequest
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	if err = c.BindJSON(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	err = services.GetServices().MessageService.PushMessage(c.GetInt64(common.UserIdKey), req)
}
