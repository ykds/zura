package entity

import (
	"github.com/ykds/zura/pkg/db"
	"github.com/ykds/zura/pkg/errors"
)

type Message struct {
	BaseModel
	FromUserId int64  `json:"from_user_id" gorm:"index"`
	ToUserId   int64  `json:"to_user_id" gorm:"index"`
	Body       string `json:"body"`
	Timestamp  int64  `json:"timestamp" gorm:"index"`
}

func (m Message) TableName() string {
	return "zura_message"
}

type GroupMessage struct {
	BaseModel
	GroupId   int64  `json:"group_id" gorm:"index"`
	UserId    int64  `json:"user_id"`
	Body      string `json:"body"`
	Timestamp int64  `json:"timestamp" gorm:"index"`
}

func (m GroupMessage) TableName() string {
	return "zura_group_message"
}

func NewMessageEntity(db *db.Database) MessageEntity {
	return &messageEntity{baseEntity{db: db}}
}

type MessageEntity interface {
	Transaction
	CreateMessage(m Message) error
	LiseNewMessage(fromUserId, toUserId, timestamp int64) ([]Message, error)
	ListHistoryMessage(fromUserId, toUserId, timestamp int64, limit int) ([]Message, error)

	CreateGroupMessage(m GroupMessage) error
	ListNewGroupMessage(groupId, timestamp int64) ([]GroupMessage, error)
	ListHistoryGroupMessage(groupId, timestamp int64, limit int) ([]GroupMessage, error)
}

type messageEntity struct {
	baseEntity
}

func (m2 messageEntity) CreateMessage(m Message) error {
	err := m2.db.Create(&m).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (m2 messageEntity) LiseNewMessage(fromUserId, toUserId, timestamp int64) ([]Message, error) {
	msg := make([]Message, 0)
	err := m2.db.Where("from_user_id = ? AND to_user_id = ? AND timestamp > ?", fromUserId, toUserId, timestamp).Order("timestamp asc").Find(&msg).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return msg, err
}

func (m2 messageEntity) ListHistoryMessage(fromUserId, toUserId, timestamp int64, limit int) ([]Message, error) {
	msg := make([]Message, 0)
	sql := m2.db.Where("((from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)) AND timestamp < ?", fromUserId, toUserId, toUserId, fromUserId, timestamp).Order("timestamp asc")
	if limit > 0 {
		sql = sql.Limit(limit)
	}
	err := sql.Find(&msg).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return msg, err
}

func (m2 messageEntity) CreateGroupMessage(m GroupMessage) error {
	err := m2.db.Create(&m).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (m2 messageEntity) ListNewGroupMessage(groupId, timestamp int64) ([]GroupMessage, error) {
	msg := make([]GroupMessage, 0)
	err := m2.db.Where("group_id = ? AND timestamp > ?", groupId, timestamp).Order("timestamp asc").Find(&msg).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return msg, err
}

func (m2 messageEntity) ListHistoryGroupMessage(groupId, timestamp int64, limit int) ([]GroupMessage, error) {
	msg := make([]GroupMessage, 0)
	sql := m2.db.Where("group_id = ? AND timestamp < ?", groupId, timestamp).Order("timestamp asc")
	if limit > 0 {
		sql = sql.Limit(limit)
	}
	err := sql.Find(&msg).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return msg, err
}
