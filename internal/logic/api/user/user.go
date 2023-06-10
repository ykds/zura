package user

import (
	"mime/multipart"
	"net/url"
	"zura/internal/logic/codec"
	"zura/internal/logic/common"
	"zura/internal/logic/entity"
	"zura/internal/logic/services"
	"zura/internal/logic/services/user"
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
	userId := c.GetInt64(common.UserIdKey)
	resp, err = services.GetServices().UserService.GetUserInfo(userId)
	avatarUrl := url.URL{
		Scheme: "http",
		Host:   c.Request.Host,
		Path:   common.StaticPath + resp.Avatar,
	}
	resp.Avatar = avatarUrl.String()
}

func UploadUserAvatar(c *gin.Context) {
	var (
		err  error
		file *multipart.FileHeader
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	file, err = c.FormFile("file")
	if err != nil {
		return
	}
	err = c.SaveUploadedFile(file, common.StaticDir+file.Filename)
	if err != nil {
		return
	}
	userId := c.GetInt64(common.UserIdKey)
	err = services.GetServices().UserService.UpdateUserInfo(userId, user.UpdateUserInfoRequest{Avatar: file.Filename})
}
