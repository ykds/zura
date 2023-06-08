package entity

import (
	"zura/pkg/db"

	"gorm.io/gorm"
)

type Friends struct {
	ID       int64 `json:"id" gorm:"primaryKey"`
	UserId   int64 `json:"user_id"`
	FriendId int64 `json:"friend_id"`
	BaseModel
}

func (f Friends) TableName() string {
	return "zura_friends"
}

func NewFriendEntity(db *db.Database) FriendEntity {
	return &friendEntity{
		db: db,
	}
}

type FriendEntity interface {
	AddFriend(userId int64, friendId int64) error
	DeleteFriend(userId int64, friendId int64) error
	ListFriend(userId int64) ([]Friends, error)
	IsFriend(userId int64, friendId int64) (bool, error)
}

type friendEntity struct {
	db *db.Database
}

func (f *friendEntity) AddFriend(userId int64, friendId int64) error {
	return f.db.Create(&Friends{UserId: userId, FriendId: friendId}).Error
}

func (f *friendEntity) DeleteFriend(userId int64, friendId int64) error {
	return f.db.Delete(&Friends{}, "user_id=? AND friend_id=?", userId, friendId).Error
}

func (f *friendEntity) ListFriend(userId int64) ([]Friends, error) {
	friends := make([]Friends, 0)
	err := f.db.Find(&friends, "user_id=?", userId).Error
	return friends, err
}

func (f *friendEntity) IsFriend(userId int64, friendId int64) (bool, error) {
	fr := Friends{}
	// TODO 测试不调用 First 是否可以
	err := f.db.Where("user_id=? AND friend_id=?", userId, friendId).Select(1).First(&fr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		return false, err
	}
	return true, nil
}
