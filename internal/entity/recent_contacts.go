package entity

import (
	"zura/pkg/db"
)

type RecentContacts struct {
	ID       int64 `json:"id" gorm:"primaryKey"`
	UserId   int64 `json:"user_id"`
	FriendId int64 `json:"friend_id"`
	IsSticky bool  `json:"is_sticky"`
	BaseModel
}

func (r RecentContacts) TableName() string {
	return "zura_recent_contacts"
}

func NewRecentContactsEntity(db *db.Database) RecentContactsEntity {
	return &recentContactsEntity{
		db: db,
	}
}

type RecentContactsEntity interface {
	AddContact(userId int64, friendId int64) error
	DeleteContact(userId int64, friendId int64) error
	ListContact(userId int64) ([]RecentContacts, error)
	UpdateContact(contactId uint, c RecentContacts) error
}

type recentContactsEntity struct {
	db *db.Database
}

func (r *recentContactsEntity) AddContact(userId int64, friendId int64) error {
	return r.db.Create(&RecentContacts{UserId: userId, FriendId: friendId}).Error
}

func (r *recentContactsEntity) DeleteContact(userId int64, friendId int64) error {
	return r.db.Delete(&RecentContacts{}, "user_id=? AND friend_id=?", userId, friendId).Error
}

func (r *recentContactsEntity) ListContact(userId int64) ([]RecentContacts, error) {
	rcList := make([]RecentContacts, 0)
	err := r.db.Find(&rcList, "user_id=?", userId).Error
	return rcList, err
}

func (r *recentContactsEntity) UpdateContact(contactId uint, c RecentContacts) error {
	return r.db.Where("id=?", contactId).Omit("user_id", "friend_id").Updates(&c).Error
}
