package session

import (
	"github.com/ykds/zura/internal/logic/codec"
	"github.com/ykds/zura/internal/logic/entity"
	"github.com/ykds/zura/pkg/errors"
	"gorm.io/gorm"
)

func NewSessionService(sessionEntity entity.SessionEntity, groupEntity entity.GroupEntity, userEntity entity.UserEntity) SessionService {
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

type SessionInfo struct {
	SessionId     int64  `json:"session_id"`
	SessionName   string `json:"session_name"`
	SessionAvatar string `json:"session_avatar"`
	IsSticky      bool   `json:"is_sticky"`
}

type UpdateUserSessionRequest struct {
	IsSticky bool `json:"is_sticky"`
}

type SessionService interface {
	CreateSession(userId int64, req CreateSessionRequest) (SessionInfo, error)
	ListSession(userId int64) ([]SessionInfo, error)
	DeleteUserSession(userId int64, id int64) error
	UpdateUserSession(userId int64, id int64, req UpdateUserSessionRequest) error
}

type sessionService struct {
	userEntity    entity.UserEntity
	sessionEntity entity.SessionEntity
	groupEntity   entity.GroupEntity
}

func (s sessionService) CreateSession(userId int64, req CreateSessionRequest) (SessionInfo, error) {
	if userId == req.TargetId {
		return SessionInfo{}, errors.WithStackByCode(codec.OpenWithSelfErrCode)
	}
	session, err := s.sessionEntity.GetUserSession(map[string]interface{}{"session_type": req.SessionType, "user_id": userId, "target_id": req.TargetId})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return SessionInfo{}, err
		}
		err = nil
		tx := s.sessionEntity.Begin()
		err = s.sessionEntity.CreateUserSessionTx(tx, entity.UserSession{SessionType: req.SessionType, UserId: userId, TargetId: req.TargetId})
		if err != nil {
			tx.Rollback()
			return SessionInfo{}, err
		}
		err = s.sessionEntity.CreateUserSessionTx(tx, entity.UserSession{SessionType: req.SessionType, UserId: req.TargetId, TargetId: userId})
		if err != nil {
			tx.Rollback()
			return SessionInfo{}, err
		}
		tx.Commit()
		session, _ = s.sessionEntity.GetUserSession(map[string]interface{}{"session_type": req.SessionType, "user_id": userId, "target_id": req.TargetId})
	}
	info := SessionInfo{
		SessionId: session.ID,
		IsSticky:  session.IsSticky,
	}
	switch req.SessionType {
	case entity.PointSession:
		user, err := s.userEntity.GetUserById(req.TargetId)
		if err != nil {
			return SessionInfo{}, err
		}
		info.SessionName = user.Username
		info.SessionAvatar = user.Avatar
	case entity.GroupSession:
		group, err := s.groupEntity.GetGroup(req.TargetId)
		if err != nil {
			return SessionInfo{}, err
		}
		info.SessionName = group.Name
		info.SessionAvatar = group.Avatar
	default:
		return SessionInfo{}, errors.WithStackByCode(codec.UnSupportSessionType)
	}
	return info, nil
}

func (s sessionService) ListSession(userId int64) ([]SessionInfo, error) {
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

	result := make([]SessionInfo, 0)
	for _, item := range session {
		info := SessionInfo{
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
	session, err := s.sessionEntity.GetUserSessionById(id)
	if err != nil {
		return err
	}
	if session.UserId != userId {
		return errors.WithStackByCode(codec.NotPermitCode)
	}
	return s.sessionEntity.DeleteUserSession(id)
}

func (s sessionService) UpdateUserSession(userId int64, id int64, req UpdateUserSessionRequest) error {
	session, err := s.sessionEntity.GetUserSessionById(id)
	if err != nil {
		return err
	}
	if session.UserId != userId {
		return errors.WithStackByCode(codec.NotPermitCode)
	}
	return s.sessionEntity.UpdateUserSession(id, entity.UserSession{
		IsSticky: req.IsSticky,
	})
}
