package entity

import (
	"context"
	"fmt"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/pkg/cache"
	"github.com/ykds/zura/pkg/db"
	"github.com/ykds/zura/pkg/errors"
	"time"
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

func NewMessageEntity(db *db.Database, cache cache.Cache) MessageEntity {
	return &messageEntity{
		baseEntity{db: db},
		cache,
	}
}

type MessageEntity interface {
	Transaction
	CreateMessage(m *Message) error
	LiseNewMessage(fromUserId, toUserId, timestamp int64) ([]Message, error)
	ListHistoryMessage(fromUserId, toUserId, timestamp int64, limit int) ([]Message, error)

	CreateGroupMessage(m *GroupMessage) error
	ListNewGroupMessage(userId, groupId, timestamp int64) ([]GroupMessage, error)
	ListHistoryGroupMessage(groupId, timestamp int64, limit int) ([]GroupMessage, error)
}

type messageEntity struct {
	baseEntity
	cache cache.Cache
}

func (m2 messageEntity) CreateMessage(m *Message) error {
	// 用于重发检查，避免插入相同的消息
	_, err := m2.cache.Get(context.Background(), fmt.Sprintf(common.MessageCacheKey, m.Timestamp))
	if err != nil {
		if errors.Is(err, cache.NotFoundErr) {
			err = m2.db.Create(&m).Error
			if err != nil {
				return errors.WithStack(err)
			}
			_ = m2.cache.Set(context.Background(), fmt.Sprintf(common.MessageCacheKey, m.Timestamp), "", 2*time.Minute)
			return nil
		}
		return err
	}
	return nil
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

func (m2 messageEntity) CreateGroupMessage(m *GroupMessage) error {
	_, err := m2.cache.Get(context.Background(), fmt.Sprintf(common.GroupMessageCacheKey, m.Timestamp))
	if err != nil {
		if errors.Is(err, cache.NotFoundErr) {
			err = m2.db.Create(&m).Error
			if err != nil {
				return errors.WithStack(err)
			}
			_ = m2.cache.Set(context.Background(), fmt.Sprintf(common.MessageCacheKey, m.Timestamp), "", 2*time.Minute)
			return nil
		}
		return err
	}
	return nil
}

func (m2 messageEntity) ListNewGroupMessage(userId, groupId, timestamp int64) ([]GroupMessage, error) {
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
