package session

import (
	"strconv"
	"zura/internal/logic/common"
	"zura/internal/logic/services"
	"zura/internal/logic/services/session"
	"zura/pkg/errors"
	"zura/pkg/response"

	"github.com/gin-gonic/gin"
)

func OpenSession(c *gin.Context) {
	var (
		err error
		req struct {
			FriendId int64 `json:"friend_id"`
		}
		resp session.SessionInfo
	)
	defer func() {
		response.HttpResponse(c, err, resp)
	}()
	if err = c.BindJSON(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	resp, err = services.GetServices().SessionService.OpenSession(c.GetInt64(common.UserIdKey), req.FriendId)
}

func ListSession(c *gin.Context) {
	var (
		err  error
		resp struct {
			Data []session.SessionInfo `json:"data"`
		}
	)
	defer func() {
		if len(resp.Data) == 0 {
			resp.Data = []session.SessionInfo{}
		}
		response.HttpResponse(c, err, resp)
	}()
	resp.Data, err = services.GetServices().SessionService.ListSession(c.GetInt64(common.UserIdKey))
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
	err = services.GetServices().SessionService.DeleteSession(c.GetInt64(common.UserIdKey), sessionId)
}
