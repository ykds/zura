package entity

import (
	"zura/pkg/db"
	"zura/pkg/errors"

	"gorm.io/gorm"
)

const (
	Normal = iota + 1
	DeletedByOne
	DeletedByTwo
	DeletedByEach
	BlackListByOne
	BlackListByTwo
	BlackListByEach
)

type Friends struct {
	User1Id int64 `json:"user1_id" gorm:"uniqueIndex:ids"`
	User2Id int64 `json:"user2_id" gorm:"uniqueIndex:ids"`
	Status  int8  `json:"status"`
	BaseModel
}

func (f Friends) TableName() string {
	return "zura_friends"
}

func NewFriendEntity(db *db.Database) FriendEntity {
	return &friendEntity{
		baseEntity{db: db},
	}
}

type FriendEntity interface {
	Transaction
	AddFriend(user1Id int64, user2Id int64) error
	AddFriendTx(tx *gorm.DB, user1Id int64, user2Id int64) error
	UpdateStatus(id int64, status int8) error
	IsFriend(user1Id int64, user2Id int64) (bool, error)
	ListFriends(userId int64) ([]Friends, error)
	GetFriend(user1Id int64, user2Id int64) (Friends, error)
}

type friendEntity struct {
	baseEntity
}

func (f *friendEntity) GetFriend(user1Id int64, user2Id int64) (Friends, error) {
	friend := Friends{}
	err := f.db.Where("(user1_id=? AND user2_id=?) OR (user2_id=? AND user1_id=?)", user1Id, user2Id, user1Id, user2Id).First(&friend).Error
	if err != nil {
		err = errors.Wrap(err, "查询好友关系失败")
	}
	return friend, err
}

func (f *friendEntity) AddFriend(user1Id int64, user2Id int64) error {
	return f.AddFriendTx(f.db.DB, user1Id, user2Id)
}

func (f *friendEntity) AddFriendTx(tx *gorm.DB, user1Id int64, user2Id int64) error {
	err := tx.Create(&Friends{User1Id: user1Id, User2Id: user2Id, Status: Normal}).Error
	if err != nil {
		err = errors.Wrap(err, "添加好友失败")
	}
	return err
}

func (f *friendEntity) ListFriends(userId int64) ([]Friends, error) {
	friends := make([]Friends, 0)
	err := f.db.Raw(
		`SELECT * from zura_friends where user1_id=? AND status IN ? 
		 UNION ALL 
		 SELECT * from zura_friends where user2_id=? AND status IN ?`, userId, []int8{Normal, DeletedByTwo}, userId, []int8{Normal, DeletedByOne}).Scan(&friends).Error
	if err != nil {
		err = errors.Wrap(err, "查询好友列表失败")
	}
	return friends, err
}

func (f *friendEntity) IsFriend(user1Id int64, user2Id int64) (bool, error) {
	fr := Friends{}
	// TODO 测试不调用 First 是否可以
	err := f.db.Where("((user1_id=? AND user2_id=?) OR (user2_id=? AND user1_id=?)) AND status=?", user1Id, user2Id, user1Id, user2Id, Normal).Select(1).First(&fr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, errors.Wrap(err, "查询好友关系失败")
	}
	return true, nil
}

func (f *friendEntity) UpdateStatus(id int64, status int8) error {
	err := f.db.Model(&Friends{}).Where("id=?", id).Update("status", status).Error
	if err != nil {
		err = errors.Wrap(err, "更新好友关系失败")
	}
	return err
}
