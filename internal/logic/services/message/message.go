package message

import (
	"encoding/json"
	"fmt"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/internal/logic/codec"
	"github.com/ykds/zura/internal/logic/entity"
	"github.com/ykds/zura/pkg/cache"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/kafka"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/proto/logic"
	"github.com/ykds/zura/proto/protocol"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

type PushMessageRequest struct {
	SessionKey int64  `json:"session_key"`
	Content    string `json:"content"`
	Timestamp  int64  `json:"timestamp"`
}

type ListMessageRequest struct {
	SessionKey int64 `form:"session_key"`
	Timestamp  int64 `form:"timestamp"`
	Limit      int   `form:"limit"`
}

type Item struct {
	ID         int64  `json:"id"`
	UniKey     int64  `json:"uni_key"`
	SessionId  int64  `json:"session_id"`
	SendUserId int64  `json:"send_user_id"`
	Body       string `json:"body"`
	Timestamp  int64  `json:"timestamp"`
}

func NewMessageService(cache cache.Cache, kafkaProducer *kafka.Producer, messageEntity entity.MessageEntity, sessionEntity entity.SessionEntity,
	groupEntity entity.GroupEntity, friendEntity entity.FriendEntity) Service {
	return &messageService{
		cache:         cache,
		kafkaProducer: kafkaProducer,
		messageEntity: messageEntity,
		sessionEntity: sessionEntity,
		friendEntity:  friendEntity,
		groupEntity:   groupEntity,
	}
}

type Service interface {
	PushMessage(userId int64, req PushMessageRequest) error
	ListNewMessage(userId int64, req ListMessageRequest) ([]Item, error)
	ListHistoryMessage(userId int64, req ListMessageRequest) ([]Item, error)
}

type messageService struct {
	cache         cache.Cache
	kafkaProducer *kafka.Producer
	messageEntity entity.MessageEntity
	sessionEntity entity.SessionEntity
	friendEntity  entity.FriendEntity
	groupEntity   entity.GroupEntity
}

func (m messageService) ListHistoryMessage(userId int64, req ListMessageRequest) ([]Item, error) {
	session, err := m.sessionEntity.GetUserSessionById(userId, req.SessionKey)
	if err != nil {
		return nil, err
	}
	result := make([]Item, 0)
	switch session.SessionType {
	case entity.PointSession:
		ok, err := m.friendEntity.IsFriend(userId, session.TargetId)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.WithStackByCode(codec.NotFriendCode)
		}
		message, err := m.messageEntity.ListHistoryMessage(userId, session.TargetId, req.Timestamp, req.Limit)
		if err != nil {
			return nil, err
		}
		for _, item := range message {
			result = append(result, Item{
				ID:         item.ID,
				UniKey:     item.Timestamp,
				SessionId:  session.ID,
				SendUserId: item.FromUserId,
				Body:       item.Body,
				Timestamp:  item.Timestamp,
			})
		}
		return result, nil
	case entity.GroupSession:
		ok, err := m.groupEntity.IsGroupMember(session.TargetId, userId)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.WithStackByCode(codec.NotGroupMember)
		}
		message, err := m.messageEntity.ListHistoryGroupMessage(session.TargetId, req.Timestamp, req.Limit)
		if err != nil {
			return nil, err
		}
		for _, item := range message {
			result = append(result, Item{
				ID:         item.ID,
				SessionId:  session.ID,
				SendUserId: item.UserId,
				Body:       item.Body,
				Timestamp:  item.Timestamp,
			})
		}
		return result, nil
	default:
		return nil, errors.WithStackByCode(codec.UnSupportSessionType)
	}
}

func (m messageService) PushMessage(userId int64, req PushMessageRequest) error {
	if time.Now().UnixMilli()-req.Timestamp > 120000 {
		return errors.WithStackByCode(codec.IllegalMsgTsCode)
	}
	session, err := m.sessionEntity.GetUserSessionById(userId, req.SessionKey)
	if err != nil {
		return err
	}
	var (
		msgId    int64
		notiUser []int64
	)
	switch session.SessionType {
	case entity.PointSession:
		ok, err := m.friendEntity.IsFriend(userId, session.TargetId)
		if err != nil {
			return err
		}
		if !ok {
			return errors.WithStackByCode(codec.NotFriendCode)
		}
		msg := &entity.Message{
			FromUserId: userId,
			ToUserId:   session.TargetId,
			Body:       req.Content,
			Timestamp:  req.Timestamp,
		}
		err = m.messageEntity.CreateMessage(msg)
		if err != nil {
			return err
		}
		msgId = msg.ID
		notiUser = []int64{session.TargetId}
	case entity.GroupSession:
		ok, err := m.groupEntity.IsGroupMember(session.TargetId, userId)
		if err != nil {
			return err
		}
		if !ok {
			return errors.WithStackByCode(codec.NotGroupMember)
		}
		msg := &entity.GroupMessage{
			GroupId:   session.TargetId,
			UserId:    userId,
			Body:      req.Content,
			Timestamp: req.Timestamp,
		}
		err = m.messageEntity.CreateGroupMessage(msg)
		if err != nil {
			return err
		}
		members, err := m.groupEntity.ListGroupMembers(session.TargetId)
		if err != nil {
			return err
		}
		msgId = msg.ID
		for _, item := range members {
			notiUser = append(notiUser, item.UserId)
		}
	default:
		return errors.WithStackByCode(codec.UnSupportSessionType)
	}

	serverUserMap := make(map[int32][]int64)
	keys := make([]string, 0)
	for _, v := range notiUser {
		keys = append(keys, fmt.Sprintf(common.UserOnlineCacheKey, v))
	}
	result, err := m.cache.MGet(context.Background(), keys...)
	if err != nil {
		return err
	}
	for i := range keys {
		if result[i] != "" {
			server, _ := strconv.ParseInt(result[i], 10, 64)
			serverUserMap[int32(server)] = append(serverUserMap[int32(server)], notiUser[i])
		}
	}
	for k, v := range serverUserMap {
		body := &logic.PushMsg{
			Op:       protocol.OpNewMsg,
			Server:   k,
			ToUserId: v,
			Message: &protocol.Message{
				Id:         msgId,
				Timestamp:  req.Timestamp,
				FromUserId: userId,
				Content:    req.Content,
			},
		}
		marshal, _ := json.Marshal(body)
		err := m.kafkaProducer.WriteMessage(context.TODO(), strconv.FormatInt(session.SessionKey, 10), marshal)
		if err != nil {
			log.Errorf("push msg error: %v", err)
		}
	}
	return nil
}

type PushMessageBody struct {
	ToUserId []int64                `json:"to_user_id"`
	Op       int32                  `json:"op"`
	Body     map[string]interface{} `json:"body"`
}

func (m messageService) ListNewMessage(userId int64, req ListMessageRequest) ([]Item, error) {
	session, err := m.sessionEntity.GetUserSessionById(userId, req.SessionKey)
	if err != nil {
		return nil, err
	}
	result := make([]Item, 0)
	switch session.SessionType {
	case entity.PointSession:
		ok, err := m.friendEntity.IsFriend(userId, session.TargetId)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.WithStackByCode(codec.NotFriendCode)
		}
		message, err := m.messageEntity.LiseNewMessage(session.TargetId, userId, req.Timestamp)
		if err != nil {
			return nil, err
		}
		for _, item := range message {
			result = append(result, Item{
				ID:         item.ID,
				UniKey:     item.Timestamp,
				SessionId:  session.ID,
				SendUserId: item.FromUserId,
				Body:       item.Body,
				Timestamp:  item.Timestamp,
			})
		}
		return result, nil
	case entity.GroupSession:
		ok, err := m.groupEntity.IsGroupMember(session.TargetId, userId)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.WithStackByCode(codec.NotGroupMember)
		}
		message, err := m.messageEntity.ListNewGroupMessage(userId, session.TargetId, req.Timestamp)
		if err != nil {
			return nil, err
		}
		for _, item := range message {
			result = append(result, Item{
				ID:         item.ID,
				SessionId:  session.ID,
				SendUserId: item.UserId,
				Body:       item.Body,
				Timestamp:  item.Timestamp,
			})
		}
		return result, nil
	default:
		return nil, errors.WithStackByCode(codec.UnSupportSessionType)
	}
}
