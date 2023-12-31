package session

import (
	"github.com/ykds/zura/internal/logic/codec"
	"github.com/ykds/zura/internal/logic/entity"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/snowflake"
	"gorm.io/gorm"
)

func NewSessionService(sessionEntity entity.SessionEntity, groupEntity entity.GroupEntity, userEntity entity.UserEntity) Service {
	return &sessionService{
		userEntity:    userEntity,
		sessionEntity: sessionEntity,
		groupEntity:   groupEntity,
	}
}

type CreateSessionRequest struct {
	SessionType int8  `json:"session_type"`
	TargetId    int64 `json:"target_id"`
}

type Info struct {
	SessionId     int64  `json:"session_id"`
	SessionKey    int64  `json:"session_key"`
	SessionName   string `json:"session_name"`
	SessionAvatar string `json:"session_avatar"`
	IsSticky      bool   `json:"is_sticky"`
}

type UpdateUserSessionRequest struct {
	IsSticky bool `json:"is_sticky"`
}

type Service interface {
	CreateSession(userId int64, req CreateSessionRequest) (Info, error)
	ListSession(userId int64) ([]Info, error)
	DeleteUserSession(userId int64, id int64) error
	UpdateUserSession(userId int64, id int64, req UpdateUserSessionRequest) error
}

type sessionService struct {
	userEntity    entity.UserEntity
	sessionEntity entity.SessionEntity
	groupEntity   entity.GroupEntity
}

func (s sessionService) CreateSession(userId int64, req CreateSessionRequest) (Info, error) {
	if userId == req.TargetId {
		return Info{}, errors.WithStackByCode(codec.OpenWithSelfErrCode)
	}
	session, err := s.sessionEntity.GetUserSession(map[string]interface{}{"session_type": req.SessionType, "user_id": userId, "target_id": req.TargetId})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return Info{}, err
		}
		tx := s.sessionEntity.Begin()

		var sessionKey int64
		userSession, err := s.sessionEntity.GetUserSession(map[string]interface{}{"session_type": req.SessionType, "user_id": req.TargetId, "target_id": userId})
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return Info{}, err
			}
			sessionKey = snowflake.NewId()
			err = s.sessionEntity.CreateUserSessionTx(tx, entity.UserSession{SessionType: req.SessionType, UserId: req.TargetId, TargetId: userId, SessionKey: sessionKey})
			if err != nil {
				tx.Rollback()
				return Info{}, err
			}
		} else {
			sessionKey = userSession.SessionKey
		}
		err = nil
		session = entity.UserSession{SessionType: req.SessionType, UserId: userId, TargetId: req.TargetId, SessionKey: sessionKey}
		err = s.sessionEntity.CreateUserSessionTx(tx, session)
		if err != nil {
			tx.Rollback()
			return Info{}, err
		}
		tx.Commit()
	}
	info := Info{
		SessionId:  session.ID,
		SessionKey: session.SessionKey,
		IsSticky:   session.IsSticky,
	}
	switch req.SessionType {
	case entity.PointSession:
		user, err := s.userEntity.GetUserById(req.TargetId)
		if err != nil {
			return Info{}, err
		}
		info.SessionName = user.Username
		info.SessionAvatar = user.Avatar
	case entity.GroupSession:
		group, err := s.groupEntity.GetGroup(req.TargetId)
		if err != nil {
			return Info{}, err
		}
		info.SessionName = group.Name
		info.SessionAvatar = group.Avatar
	default:
		return Info{}, errors.WithStackByCode(codec.UnSupportSessionType)
	}
	return info, nil
}

func (s sessionService) ListSession(userId int64) ([]Info, error) {
	session, err := s.sessionEntity.ListSession(userId)
	if err != nil {
		return nil, err
	}
	userIds := make([]int64, 0)
	groupIds := make([]int64, 0)
	for _, item := range session {
		switch item.SessionType {
		case entity.PointSession:
			userIds = append(userIds, item.TargetId)
		case entity.GroupSession:
			groupIds = append(groupIds, item.TargetId)
		}
	}
	users, err := s.userEntity.ListUserById(userIds)
	if err != nil {
		return nil, err
	}
	usersMap := make(map[int64]entity.User)
	for _, u := range users {
		usersMap[u.ID] = u
	}
	groups, err := s.groupEntity.ListGroupById(groupIds)
	if err != nil {
		return nil, err
	}
	groupMap := make(map[int64]entity.Group)
	for _, g := range groups {
		groupMap[g.ID] = g
	}

	result := make([]Info, 0)
	for _, item := range session {
		info := Info{
			SessionId: item.ID,
			IsSticky:  item.IsSticky,
		}
		switch item.SessionType {
		case entity.PointSession:
			user, ok := usersMap[item.TargetId]
			if !ok {
				continue
			}
			info.SessionName = user.Username
			info.SessionAvatar = user.Avatar
		case entity.GroupSession:
			group, ok := groupMap[item.TargetId]
			if !ok {
				continue
			}
			info.SessionName = group.Name
			info.SessionAvatar = group.Avatar
		}
		result = append(result, info)
	}
	return result, nil
}

func (s sessionService) DeleteUserSession(userId int64, id int64) error {
	_, err := s.sessionEntity.GetUserSessionById(userId, id)
	if err != nil {
		return err
	}
	//if session.UserId != userId {
	//	return errors.WithStackByCode(codec.NotPermitCode)
	//}
	return s.sessionEntity.DeleteUserSession(id)
}

func (s sessionService) UpdateUserSession(userId int64, id int64, req UpdateUserSessionRequest) error {
	_, err := s.sessionEntity.GetUserSessionById(userId, id)
	if err != nil {
		return err
	}
	//if session.UserId != userId {
	//	return errors.WithStackByCode(codec.NotPermitCode)
	//}
	return s.sessionEntity.UpdateUserSession(id, entity.UserSession{
		IsSticky: req.IsSticky,
	})
}
