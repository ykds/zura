package entity

import (
	"zura/pkg/db"
)

type User struct {
	Password string `json:"password" gorm:"password"`
	Salt     string `json:"salt" gorm:"salt"`
	UserInfo
	BaseModel
}

type UserInfo struct {
	UserId   int64  `json:"user_id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"username"`
	Phone    string `json:"phone" gorm:"phone"`
	Email    string `json:"email" gorm:"email"`
	Avatar   string `json:"avatar" gorm:"avatar"`
}

func (u User) TableName() string {
	return "zura_users"
}

func NewUserEntity(db *db.Database) UserEntity {
	return &userEntity{
		db: db,
	}
}

type UserEntity interface {
	CreateUser(user User) error
	GetUserById(userId int64) (User, error)
	GetUser(where map[string]interface{}) (User, error)
	UpdateUser(userId int64, user User) error
}

type userEntity struct {
	db *db.Database
}

func (u *userEntity) CreateUser(user User) error {
	return u.db.Create(&user).Error
}

func (u *userEntity) GetUserById(userId int64) (User, error) {
	user := User{}
	err := u.db.First(&user, "user_id=?", userId).Error
	return user, err
}

func (u *userEntity) GetUser(where map[string]interface{}) (User, error) {
	user := User{}
	err := u.db.Where(where).First(&user).Error
	return user, err
}

func (u *userEntity) UpdateUser(userId int64, user User) error {
	return u.db.Where("user_id=?", userId).Updates(&user).Error
}
