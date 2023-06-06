package user

import "gorm.io/gorm"

type User struct {
	UserId   int64  `json:"user_id" gorm:"primaryKey"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	gorm.Model
}

func (u User) TableName() string {
	return "zira_users"
}

type Friends struct {
	gorm.Model
	UserId int64 `json:"user_id"`
	FriendId int64 `json:"friend_id"`
}

func (f Friends) TableName() string {
	return "zira_friends"
}

type RecentContacts struct {
	gorm.Model
	UserId int64 `json:"user_id"`
	FriendId int64 `json:"friend_id"`
}

func (r RecentContacts) TableName() string {
	return "zira_recent_contacts"
}