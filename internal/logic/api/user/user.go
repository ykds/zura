package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/internal/logic/codec"
	"github.com/ykds/zura/internal/logic/entity"
	"github.com/ykds/zura/internal/logic/services"
	"github.com/ykds/zura/internal/logic/services/user"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/response"
)

func Register(c *gin.Context) {
	var (
		err error
		req user.RegisterRequest
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	if err = c.BindJSON(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}

	if req.RegisterType == user.PhoneType && req.Phone == "" {
		err = errors.New(codec.PhoneEmptyStatus)
		return
	} else if req.RegisterType == user.EmailType && req.Email == "" {
		err = errors.New(codec.EmailEmptyStatus)
		return
	} else if req.RegisterType == user.UsernameType && req.Username == "" {
		err = errors.New(codec.UsernameEmptyStatus)
		return
	}

	if req.Password == "" || req.ConfirmPassword == "" {
		err = errors.New(errors.ParameterErrorStatus)
		return
	}

	err = services.GetServices().UserService.Register(req)
}

func Login(c *gin.Context) {
	var (
		err  error
		req  user.LoginRequest
		resp user.LoginResponse
	)
	defer func() {
		response.HttpResponse(c, err, resp)
	}()
	if err = c.BindJSON(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	if req.LoginType == user.PhoneType && req.Phone == "" {
		err = errors.New(codec.PhoneEmptyStatus)
		return
	} else if req.LoginType == user.EmailType && req.Email == "" {
		err = errors.New(codec.EmailEmptyStatus)
		return
	} else if req.LoginType == user.UsernameType && req.Username == "" {
		err = errors.New(codec.UsernameEmptyStatus)
		return
	}
	resp, err = services.GetServices().UserService.Login(req)
}

func GetUserInfo(c *gin.Context) {
	var (
		err  error
		resp entity.UserInfo
	)
	defer func() {
		response.HttpResponse(c, err, resp)
	}()
	userId := c.GetInt64(common.UserIdKey)
	resp, err = services.GetServices().UserService.GetUserInfo(userId)
	resp.Avatar = common.ParseAvatarUrl(c, resp.Avatar)
}

func UpdateInfo(c *gin.Context) {
	var (
		err error
		req user.UpdateUserInfoRequest
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	if err = c.BindJSON(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	userId := c.GetInt64(common.UserIdKey)
	err = services.GetServices().UserService.UpdateUserInfo(userId, req)
}

func SearchUser(c *gin.Context) {
	var (
		err  error
		req  user.SearchUsersRequest
		resp struct {
			Data []user.OtherUserInfo `json:"data"`
		}
	)
	defer func() {
		if len(resp.Data) == 0 {
			resp.Data = make([]user.OtherUserInfo, 0)
		}
		response.HttpResponse(c, err, resp)
	}()
	if err = c.BindQuery(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	resp.Data, err = services.GetServices().UserService.SearchUser(req)
	if err == nil {
		for i, item := range resp.Data {
			resp.Data[i].Avatar = common.ParseAvatarUrl(c, item.Avatar)
		}
	}
}
