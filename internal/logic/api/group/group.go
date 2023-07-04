package group

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/internal/logic/services"
	"github.com/ykds/zura/internal/logic/services/group"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/response"
	"strconv"
)

func CreateGroup(c *gin.Context) {
	var (
		err error
		req group.CreateGroupRequest
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	if err = c.BindJSON(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	err = services.GetServices().GroupService.CreateGroup(c.GetInt64(common.UserIdKey), req)
}

func ListGroup(c *gin.Context) {
	var (
		err  error
		resp struct {
			Data []group.Info `json:"data"`
		}
	)
	defer func() {
		if len(resp.Data) == 0 {
			resp.Data = make([]group.Info, 0)
		}
		response.HttpResponse(c, err, resp)
	}()
	resp.Data, err = services.GetServices().GroupService.ListGroup(c.GetInt64(common.UserIdKey))
	if err != nil {
		return
	}
	for i, item := range resp.Data {
		resp.Data[i].Avatar = common.ParseAvatarUrl(c, item.Avatar)
	}
}

func SearchGroup(c *gin.Context) {
	var (
		err  error
		req  group.SearchGroupRequest
		resp struct {
			Data []group.Info `json:"data"`
		}
	)
	defer func() {
		if len(resp.Data) == 0 {
			resp.Data = make([]group.Info, 0)
		}
		response.HttpResponse(c, err, resp)
	}()
	resp.Data, err = services.GetServices().GroupService.SearchGroup(req)
	if err != nil {
		return
	}
	for i, item := range resp.Data {
		resp.Data[i].Avatar = common.ParseAvatarUrl(c, item.Avatar)
	}
}

func UpdateGroup(c *gin.Context) {
	var (
		err error
		req group.UpdateGroupRequest
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	if err = c.BindJSON(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	err = services.GetServices().GroupService.UpdateGroup(c.GetInt64(common.UserIdKey), req)
}

func DismissGroup(c *gin.Context) {
	var (
		err error
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	param := c.Param("group_id")
	groupId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	err = services.GetServices().GroupService.DismissGroup(c.GetInt64(common.UserIdKey), groupId)
}

func AddGroupMember(c *gin.Context) {
	var (
		err error
		req group.AddGroupMemberRequest
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	if err = c.BindJSON(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	err = services.GetServices().GroupService.AddGroupMember(c.GetInt64(common.UserIdKey), req)
}

func RemoveGroupMember(c *gin.Context) {
	var (
		err error
		req group.RemoveGroupMemberRequest
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	if err = c.BindJSON(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	err = services.GetServices().GroupService.RemoveGroupMember(c.GetInt64(common.UserIdKey), req)
}

func UpdateMemberInfo(c *gin.Context) {
	var (
		err error
		req group.UpdateMemberInfoRequest
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	if err = c.BindJSON(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	err = services.GetServices().GroupService.UpdateMemberInfo(c.GetInt64(common.UserIdKey), req)
}

func ListGroupMembers(c *gin.Context) {
	var (
		err  error
		resp struct {
			Data []group.MemberInfo `json:"data"`
		}
	)
	defer func() {
		if len(resp.Data) == 0 {
			resp.Data = make([]group.MemberInfo, 0)
		}
		response.HttpResponse(c, err, resp)
	}()
	param := c.Param("group_id")
	groupId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	resp.Data, err = services.GetServices().GroupService.ListGroupMembers(c.GetInt64(common.UserIdKey), groupId)
	if err != nil {
		return
	}
	for i, item := range resp.Data {
		resp.Data[i].Avatar = common.ParseAvatarUrl(c, item.Avatar)
	}
}

func ChangeMemberRole(c *gin.Context) {
	var (
		err error
		req group.ChangeMemberRoleRequest
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	if err = c.BindJSON(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	err = services.GetServices().GroupService.ChangeMemberRole(c.GetInt64(common.UserIdKey), req)
}
