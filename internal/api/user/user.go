package user

import (
	"zura/internal/codec"
	"zura/internal/entity"
	"zura/internal/services"
	"zura/internal/services/user"
	"zura/pkg/errors"
	"zura/pkg/response"

	"github.com/gin-gonic/gin"
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
	userId := c.GetInt64("userId")
	resp, err = services.GetServices().UserService.GetUserInfo(userId)
}
