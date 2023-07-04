package entity

import (
	"github.com/ykds/zura/pkg/db"
	"github.com/ykds/zura/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Password string `json:"password" gorm:"password"`
	Salt     string `json:"salt" gorm:"salt"`
	UserInfo
}

type UserInfo struct {
	BaseModel
	Avatar            string    `json:"avatar" gorm:"avatar"`
	Username          string    `json:"username" gorm:"username"`
	Phone             string    `json:"phone" gorm:"phone"`
	Email             string    `json:"email" gorm:"email"`
	UpdatedPhoneAt    time.Time `json:"updated_phone_at"`
	UpdatedEmailAt    time.Time `json:"updated_email_at"`
	UpdatedUsernameAt time.Time `json:"updated_username_at"`
}

func (u *User) TableName() string {
	return "zura_users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	err = u.BaseModel.BeforeCreate(tx)
	if err != nil {
		return err
	}
	u.UpdatedPhoneAt = time.Now()
	u.UpdatedEmailAt = time.Now()
	u.UpdatedUsernameAt = time.Now()
	return
}

func (u *User) BeforeUpdate(_ *gorm.DB) (err error) {
	if u.Phone != "" {
		u.UpdatedPhoneAt = time.Now()
	}
	if u.Email != "" {
		u.UpdatedEmailAt = time.Now()
	}
	if u.Username != "" {
		u.UpdatedUsernameAt = time.Now()
	}
	return
}

func NewUserEntity(db *db.Database) UserEntity {
	return &userEntity{
		db: db,
	}
}

type UserEntity interface {
	CreateUser(user User) error
	GetUserById(userId int64) (User, error)
	ListUserById(userId []int64) ([]User, error)
	GetUser(where map[string]interface{}) (User, error)
	UpdateUser(userId int64, user User) error
	ListUser(where map[string]interface{}) ([]User, error)
}

type userEntity struct {
	db *db.Database
}

func (u *userEntity) ListUser(where map[string]interface{}) ([]User, error) {
	users := make([]User, 0)
	err := u.db.Where(where).Find(&users).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return users, err
}

func (u *userEntity) CreateUser(user User) error {
	err := u.db.Create(&user).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (u *userEntity) GetUserById(userId int64) (User, error) {
	user := User{}
	err := u.db.First(&user, "id=?", userId).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return user, err
}

func (u *userEntity) ListUserById(userId []int64) ([]User, error) {
	user := make([]User, 0)
	err := u.db.Find(&user, "id IN ?", userId).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return user, err
}

func (u *userEntity) GetUser(where map[string]interface{}) (User, error) {
	user := User{}
	err := u.db.Where(where).First(&user).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return user, err
}

func (u *userEntity) UpdateUser(userId int64, user User) error {
	err := u.db.Where("id=?", userId).Updates(&user).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}
