package friend_application

import (
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/internal/logic/codec"
	"github.com/ykds/zura/internal/logic/services"
	"github.com/ykds/zura/internal/logic/services/friend_application"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Apply(c *gin.Context) {
	var (
		err error
		req friend_application.ApplyRequest
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	if err = c.BindJSON(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	if req.UserId == 0 {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	userId := c.GetInt64(common.UserIdKey)
	if req.UserId == userId {
		err = errors.New(codec.ApplyMySelfErrorCode)
		return
	}
	err = services.GetServices().FriendApplicationService.ApplyFriend(userId, req)
}

func ListApplications(c *gin.Context) {
	var (
		err  error
		resp struct {
			Data []friend_application.Application `json:"data"`
		}
	)
	defer func() {
		if len(resp.Data) == 0 {
			resp.Data = []friend_application.Application{}
		}
		response.HttpResponse(c, err, resp)
	}()
	userId := c.GetInt64(common.UserIdKey)
	resp.Data, err = services.GetServices().FriendApplicationService.ListApplications(userId)
}

func UpdateApplicationStatus(c *gin.Context) {
	var (
		err error
		req struct {
			Status int8 `json:"status"`
		}
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	if err = c.BindJSON(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	s := c.Param("id")
	if s == "" {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	var id int64
	id, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}
	err = services.GetServices().FriendApplicationService.UpdateApplicationStatus(c.GetInt64(common.UserIdKey), id, req.Status)
}

func DeleteApplication(c *gin.Context) {
	var (
		err error
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	userId := c.GetInt64(common.UserIdKey)
	s := c.Param("id")
	if s == "" {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	var id int64
	id, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}
	err = services.GetServices().FriendApplicationService.DeleteApplication(id, userId)
}
