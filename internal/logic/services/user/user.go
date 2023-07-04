package user

import (
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/internal/logic/codec"
	"github.com/ykds/zura/internal/logic/config"
	"github.com/ykds/zura/internal/logic/entity"
	"github.com/ykds/zura/internal/logic/services/verify_code"
	"github.com/ykds/zura/pkg/cache"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/random"
	"github.com/ykds/zura/pkg/token"
	"gorm.io/gorm"
	"time"
)

const (
	UsernameType = "username"
	PhoneType    = "phone"
	EmailType    = "email"
)

func NewUserService(cache cache.Cache, userEntity entity.UserEntity, verifyCodeService verify_code.VerifyCodeService) Service {
	return &userService{
		cache:             cache,
		userEntity:        userEntity,
		verifyCodeService: verifyCodeService,
	}
}

type RegisterRequest struct {
	RegisterType    string `json:"register_type"`
	Username        string `json:"username"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
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
	Username   string `json:"username"`
	Avatar     string `json:"avatar"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	VerifyCode string `json:"verify_code"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type SearchUsersRequest struct {
	SearchType string `form:"search_type"`
	Username   string `form:"username"`
	Phone      string `form:"phone"`
	Email      string `form:"email"`
}

type OtherUserInfo struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type Service interface {
	Register(req RegisterRequest) error
	Login(req LoginRequest) (LoginResponse, error)
	Logout(userId int64) error
	GetUserInfo(userId int64) (entity.UserInfo, error)
	UpdateUserInfo(userId int64, req UpdateUserInfoRequest) error
	ChangePassword(userId int64, req ChangePasswordRequest) error
	SearchUser(req SearchUsersRequest) ([]OtherUserInfo, error)

	Connect(ctx context.Context, userId int64, serverId int32) error
	DisConnect(ctx context.Context, userId int64) error
	HeartBeat(ctx context.Context, userId int64) error
}

type userService struct {
	cache             cache.Cache
	userEntity        entity.UserEntity
	verifyCodeService verify_code.VerifyCodeService
}

func (u *userService) Connect(ctx context.Context, userId int64, serverId int32) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return u.cache.Set(ctx, fmt.Sprintf(common.UserOnlineCacheKey, userId), serverId, time.Duration(config.GetConfig().Session.HeartbeatInterval)*time.Second)
}

func (u *userService) DisConnect(ctx context.Context, userId int64) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return u.cache.Del(ctx, fmt.Sprintf(common.UserOnlineCacheKey, userId))
}

func (u *userService) HeartBeat(ctx context.Context, userId int64) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return u.cache.Expire(ctx, fmt.Sprintf(common.UserOnlineCacheKey, userId), time.Duration(config.GetConfig().Session.HeartbeatInterval)*time.Second)
}

func (u *userService) SearchUser(req SearchUsersRequest) ([]OtherUserInfo, error) {
	where := make(map[string]interface{})
	switch req.SearchType {
	case UsernameType:
		where["username"] = req.Username
	case PhoneType:
		where["phone"] = req.Phone
	case EmailType:
		where["email"] = req.Email
	default:
		return nil, errors.WithStackByCode(codec.UnSupportedTypeStatus)
	}
	users, err := u.userEntity.ListUser(where)
	if err != nil {
		return nil, err
	}
	infos := make([]OtherUserInfo, 0)
	for _, u := range users {
		infos = append(infos, OtherUserInfo{
			UserId:   u.ID,
			Username: u.Username,
			Avatar:   u.Avatar,
		})
	}
	return infos, nil
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
		return errors.New(codec.UnSupportedTypeStatus)
	}
	_, err := u.userEntity.GetUser(where)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		} else {
			return errors.WithStack(err)
		}
	} else {
		return errors.New(codec.UserRegisteredStatus)
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
		return LoginResponse{}, errors.New(codec.UnSupportedTypeStatus)
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
	t, err := token.NewToken(info.ID)
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
	info, err := u.GetUserInfo(userId)
	if err != nil {
		return err
	}
	if req.Username != "" {
		if !info.UpdatedUsernameAt.AddDate(1, 0, 0).After(time.Now()) {
			return errors.Newf(codec.YearUpdateLimitStatus, "用户名")
		}
	}
	if req.Phone != "" {
		if !info.UpdatedPhoneAt.AddDate(1, 0, 0).After(time.Now()) {
			return errors.Newf(codec.YearUpdateLimitStatus, "手机号")
		}
		ok, err := u.verifyCodeService.VerifyCode(fmt.Sprintf(common.UpdateUserPhoneCacheKey, userId), req.VerifyCode)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New(codec.VerifyStatusWrongStatus)
		}
	}
	if req.Email != "" {
		if !info.UpdatedEmailAt.AddDate(1, 0, 0).After(time.Now()) {
			return errors.Newf(codec.YearUpdateLimitStatus, "邮箱")
		}
		ok, err := u.verifyCodeService.VerifyCode(fmt.Sprintf(common.UpdateUserEmailCacheKey, userId), req.VerifyCode)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New(codec.VerifyStatusWrongStatus)
		}
	}
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
