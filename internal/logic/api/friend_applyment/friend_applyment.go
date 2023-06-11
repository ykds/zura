package friend_applyment

import (
	"strconv"
	"zura/internal/logic/codec"
	"zura/internal/logic/common"
	"zura/internal/logic/services"
	"zura/internal/logic/services/friend_applyment"
	"zura/pkg/errors"
	"zura/pkg/response"

	"github.com/gin-gonic/gin"
)

func Apply(c *gin.Context) {
	var (
		err error
		req friend_applyment.ApplyRequest
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
	err = services.GetServices().FriendApplymentService.ApplyFriend(userId, req)
}

func ListApplyments(c *gin.Context) {
	var (
		err  error
		resp struct {
			Data []friend_applyment.Applyment `json:"data"`
		}
	)
	defer func() {
		if len(resp.Data) == 0 {
			resp.Data = []friend_applyment.Applyment{}
		}
		response.HttpResponse(c, err, resp)
	}()
	userId := c.GetInt64(common.UserIdKey)
	resp.Data, err = services.GetServices().FriendApplymentService.ListApplyments(userId)
}

func UpdateApplymentStatus(c *gin.Context) {
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
	err = services.GetServices().FriendApplymentService.UpdateApplymentStatus(c.GetInt64(common.UserIdKey), id, req.Status)
}

func DeleteApplyment(c *gin.Context) {
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
	err = services.GetServices().FriendApplymentService.DeleteApplyment(id, userId)
}
