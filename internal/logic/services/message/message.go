package message

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/internal/logic/codec"
	"github.com/ykds/zura/internal/logic/entity"
	"github.com/ykds/zura/pkg/cache"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/proto/comet"
	"time"
)

type PushMessageRequest struct {
	SessionId int64  `json:"session_id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

type ListMessageRequest struct {
	SessionId int64 `form:"session_id"`
	Timestamp int64 `form:"timestamp"`
	Limit     int   `form:"limit"`
}

type MessageItem struct {
	ID         int64  `json:"id"`
	UniKey     int64  `json:"uni_key"`
	SessionId  int64  `json:"session_id"`
	SendUserId int64  `json:"send_user_id"`
	Body       string `json:"body"`
	Timestamp  int64  `json:"timestamp"`
}

func NewMessageService(cometClient comet.CometClient, messageEntity entity.MessageEntity, sessionEntity entity.SessionEntity,
	groupEntity entity.GroupEntity, friendEntity entity.FriendEntity) MessageService {
	return &messageService{
		cometClient:   cometClient,
		messageEntity: messageEntity,
		sessionEntity: sessionEntity,
		friendEntity:  friendEntity,
		groupEntity:   groupEntity,
	}
}

type MessageService interface {
	PushMessage(userId int64, req PushMessageRequest) error
	ListNewMessage(userId int64, req ListMessageRequest) ([]MessageItem, error)
	ListHistoryMessage(userId int64, req ListMessageRequest) ([]MessageItem, error)
}

type messageService struct {
	cometClient   comet.CometClient
	messageEntity entity.MessageEntity
	sessionEntity entity.SessionEntity
	friendEntity  entity.FriendEntity
	groupEntity   entity.GroupEntity
}

func (m messageService) ListHistoryMessage(userId int64, req ListMessageRequest) ([]MessageItem, error) {
	session, err := m.sessionEntity.GetUserSessionById(req.SessionId)
	if err != nil {
		return nil, err
	}
	result := make([]MessageItem, 0)
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
			result = append(result, MessageItem{
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
			result = append(result, MessageItem{
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
	session, err := m.sessionEntity.GetUserSessionById(req.SessionId)
	if err != nil {
		return err
	}
	var notiUser []int64
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
		msgByte, _ := json.Marshal(msg)
		_ = cache.GetGloMemCache().LPush(context.Background(), fmt.Sprintf(common.UnackMessageCacheKey, msg.FromUserId, msg.ToUserId), string(msgByte), time.Minute)
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
		msgByte, _ := json.Marshal(msg)
		_ = cache.GetGloMemCache().LPush(context.Background(), fmt.Sprintf(common.UnackGroupMessageCacheKey, msg.GroupId), string(msgByte), time.Minute)
		members, err := m.groupEntity.ListGroupMembers(session.TargetId)
		if err != nil {
			return err
		}
		for _, item := range members {
			notiUser = append(notiUser, item.UserId)
		}
	default:
		return errors.WithStackByCode(codec.UnSupportSessionType)
	}
	session2, err := m.sessionEntity.GetUserSession(map[string]interface{}{"user_id": session.TargetId, "target_id": userId})
	if err != nil {
		return err
	}
	body, _ := json.Marshal(map[string]interface{}{"op": comet.Op_NewMsg, "body": map[string]interface{}{"session_id": session2.ID}})
	_, err = m.cometClient.PushNotification(context.Background(), &comet.PushNotificationRequest{
		ToUserId: notiUser,
		Body:     body,
	})
	return err
}

func (m messageService) ListNewMessage(userId int64, req ListMessageRequest) ([]MessageItem, error) {
	session, err := m.sessionEntity.GetUserSessionById(req.SessionId)
	if err != nil {
		return nil, err
	}
	result := make([]MessageItem, 0)
	switch session.SessionType {
	case entity.PointSession:
		ok, err := m.friendEntity.IsFriend(userId, session.TargetId)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, errors.WithStackByCode(codec.NotFriendCode)
		}
		var message []entity.Message
		messList, err := cache.GetGloMemCache().LRange(context.Background(), fmt.Sprintf(common.UnackMessageCacheKey, session.TargetId, userId), 0, -1)
		if err == nil {
			for _, mess := range messList {
				item := entity.Message{}
				_ = json.Unmarshal([]byte(mess), &item)
				if item.Timestamp > req.Timestamp {
					message = append(message, item)
				}
			}
			_ = cache.GetGloMemCache().LRem(context.Background(), fmt.Sprintf(common.UnackMessageCacheKey, session.TargetId, userId), 0, int64(len(message)))
		} else {
			message, err = m.messageEntity.LiseNewMessage(session.TargetId, userId, req.Timestamp)
			if err != nil {
				return nil, err
			}
		}

		for _, item := range message {
			result = append(result, MessageItem{
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

		var message []entity.GroupMessage
		messList, err := cache.GetGloMemCache().LRange(context.Background(), fmt.Sprintf(common.UnackMessageCacheKey, session.TargetId, userId), 0, -1)
		if err == nil {
			for _, mess := range messList {
				item := entity.GroupMessage{}
				_ = json.Unmarshal([]byte(mess), &item)
				if item.Timestamp > req.Timestamp {
					message = append(message, item)
				}
			}
			_ = cache.GetGloMemCache().LRem(context.Background(), fmt.Sprintf(common.UnackGroupMessageCacheKey, session.TargetId), 0, int64(len(message)))
		} else {
			message, err = m.messageEntity.ListNewGroupMessage(session.TargetId, req.Timestamp)
			if err != nil {
				return nil, err
			}
		}

		for _, item := range message {
			result = append(result, MessageItem{
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
