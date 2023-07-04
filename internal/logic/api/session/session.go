package session

import (
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/internal/logic/entity"
	"github.com/ykds/zura/internal/logic/services"
	"github.com/ykds/zura/internal/logic/services/session"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateSession(c *gin.Context) {
	var (
		err  error
		req  session.CreateSessionRequest
		resp session.Info
	)
	defer func() {
		response.HttpResponse(c, err, resp)
	}()
	if err = c.BindJSON(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	req.SessionType = entity.PointSession
	resp, err = services.GetServices().SessionService.CreateSession(c.GetInt64(common.UserIdKey), req)
}

func ListSession(c *gin.Context) {
	var (
		err  error
		resp struct {
			Data []session.Info `json:"data"`
		}
	)
	defer func() {
		if len(resp.Data) == 0 {
			resp.Data = []session.Info{}
		}
		response.HttpResponse(c, err, resp)
	}()
	resp.Data, err = services.GetServices().SessionService.ListSession(c.GetInt64(common.UserIdKey))
	if err != nil {
		return
	}
	for i, item := range resp.Data {
		resp.Data[i].SessionAvatar = common.ParseAvatarUrl(c, item.SessionAvatar)
	}
}

func DeleteSession(c *gin.Context) {
	var (
		err error
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	id := c.Param("session_id")
	var sessionId int64
	sessionId, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		return
	}
	err = services.GetServices().SessionService.DeleteUserSession(c.GetInt64(common.UserIdKey), sessionId)
}

func UpdateSession(c *gin.Context) {
	var (
		err error
		req session.UpdateUserSessionRequest
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	if err = c.BindJSON(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	id := c.Param("session_id")
	var sessionId int64
	sessionId, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		return
	}
	err = services.GetServices().SessionService.UpdateUserSession(c.GetInt64(common.UserIdKey), sessionId, req)
}
