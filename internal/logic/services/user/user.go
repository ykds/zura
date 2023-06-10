package user

import (
	"crypto/sha256"
	"fmt"
	"zura/internal/logic/codec"
	"zura/internal/logic/common"
	"zura/internal/logic/entity"
	"zura/internal/logic/services/verify_code"
	"zura/pkg/errors"
	"zura/pkg/random"
	"zura/pkg/token"

	"gorm.io/gorm"
)

const (
	UsernameType = "username"
	PhoneType    = "phone"
	EmailType    = "email"
)

func NewUserService(userEntity entity.UserEntity, verifyCodeService verify_code.VerifyCodeService) UserService {
	return &userService{
		userEntity:        userEntity,
		verifyCodeService: verifyCodeService,
	}
}

type RegisterRequest struct {
	RegisterType    string `json:"register_type"`
	Username        string `json:"username"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	Avatar          string `json:"avatar"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	VerifyCode      string `json:"verify_code"`
}

type LoginRequest struct {
	LoginType string `json:"login_type"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UpdateUserInfoRequest struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UserService interface {
	Register(req RegisterRequest) error
	Login(req LoginRequest) (LoginResponse, error)
	Logout(userId int64) error
	GetUserInfo(userId int64) (entity.UserInfo, error)
	UpdateUserInfo(userId int64, req UpdateUserInfoRequest) error
	ChangePassword(userId int64, req ChangePasswordRequest) error
}

type userService struct {
	userEntity        entity.UserEntity
	verifyCodeService verify_code.VerifyCodeService
}

func (u *userService) Register(req RegisterRequest) error {
	where := make(map[string]interface{})
	cacheKey := ""
	switch req.RegisterType {
	case PhoneType:
		where["phone"] = req.Phone
		cacheKey = fmt.Sprintf(common.RegisterVerifyCodeKey, PhoneType, req.Phone)
	case EmailType:
		where["email"] = req.Email
		cacheKey = fmt.Sprintf(common.RegisterVerifyCodeKey, EmailType, req.Email)
	case UsernameType:
		where["username"] = req.Username
		cacheKey = fmt.Sprintf(common.RegisterVerifyCodeKey, UsernameType, req.Username)
	default:
		return errors.New(codec.UnSupportedRegisterTypeStatus)
	}
	_, err := u.userEntity.GetUser(where)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		} else {
			return errors.WithStack(err)
		}
	} else {
		return errors.New(codec.UserRegisterdStatus)
	}

	if req.Password != req.ConfirmPassword {
		return errors.New(codec.PasswordNotConsistentStatus)
	}

	if req.RegisterType != UsernameType {
		ok, err := u.verifyCodeService.VerifyCode(cacheKey, req.VerifyCode)
		if err != nil {
			return errors.WithStack(err)
		}
		if !ok {
			return errors.New(codec.VerifyStatusWrongStatus)

		}
	}

	salt := random.RandStr(16)
	password := hashPassword(req.Password, salt)
	err = u.userEntity.CreateUser(entity.User{
		Salt:     salt,
		Password: password,
		UserInfo: entity.UserInfo{
			Username: req.Username,
			Phone:    req.Phone,
			Email:    req.Email,
			Avatar:   req.Avatar,
		},
	})
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (u *userService) Login(req LoginRequest) (LoginResponse, error) {
	where := make(map[string]interface{})
	switch req.LoginType {
	case PhoneType:
		where["phone"] = req.Phone
	case EmailType:
		where["email"] = req.Email
	case UsernameType:
		where["username"] = req.Username
	default:
		return LoginResponse{}, errors.New(codec.UnSupportedLoginTypeStatus)
	}
	info, err := u.userEntity.GetUser(where)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return LoginResponse{}, errors.New(codec.UserNotFoundStatus)
		}
		return LoginResponse{}, errors.WithStack(err)
	}
	if !comparePassword(req.Password, info.Password, info.Salt) {
		return LoginResponse{}, errors.New(codec.PasswordWrongStatus)
	}
	t, err := token.NewToken(info.UserId)
	if err != nil {
		return LoginResponse{}, errors.WithStack(err)
	}
	return LoginResponse{
		Token: t,
	}, nil
}

func (u *userService) Logout(userId int64) error {
	return nil
}

func (u *userService) GetUserInfo(userId int64) (entity.UserInfo, error) {
	user, err := u.userEntity.GetUserById(userId)
	if err != nil {
		return entity.UserInfo{}, errors.WithStack(err)
	}
	return user.UserInfo, nil
}

func (u *userService) UpdateUserInfo(userId int64, req UpdateUserInfoRequest) error {
	return u.userEntity.UpdateUser(userId, entity.User{
		UserInfo: entity.UserInfo{
			Username: req.Username,
			Phone:    req.Phone,
			Email:    req.Email,
			Avatar:   req.Avatar,
		},
	})
}

func (u *userService) ChangePassword(userId int64, req ChangePasswordRequest) error {
	user, err := u.userEntity.GetUserById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(codec.UserNotFoundStatus)
		}
		return errors.WithStack(err)
	}
	if req.NewPassword != req.ConfirmPassword {
		return errors.New(codec.PasswordNotConsistentStatus)
	}
	if !comparePassword(req.OldPassword, user.Password, user.Salt) {
		return errors.New(codec.OldPasswordWrongStatus)
	}
	password := hashPassword(req.NewPassword, user.Salt)
	return u.userEntity.UpdateUser(userId, entity.User{Password: password})
}

func hashPassword(password string, salt string) string {
	passwordByte := sha256.Sum256([]byte(password + salt))
	return fmt.Sprintf("%x", passwordByte)
}

func comparePassword(testpasswd string, passwd string, salt string) bool {
	return hashPassword(testpasswd, salt) == passwd
}
