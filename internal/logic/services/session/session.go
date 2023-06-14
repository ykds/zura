package session

import (
	"github.com/ykds/zura/internal/logic/codec"
	"github.com/ykds/zura/internal/logic/entity"
	"github.com/ykds/zura/internal/logic/services/friends"
	"github.com/ykds/zura/pkg/errors"
	"gorm.io/gorm"
	"strconv"
)

func NewSessionService(sessionEntity entity.SessionEntity, friendsService friends.FriendsService) SessionService {
	return &sessionService{
		sessionEntity:  sessionEntity,
		friendsService: friendsService,
	}
}

type SessionInfo struct {
	ID          int64              `json:"id"`
	SessionType int8               `json:"session_type"`
	SessionName string             `json:"session_name"`
	Members     []int64            `json:"members,omitempty"`
	FriendId    int64              `json:"friend_id,omitempty"`
	Setting     SessionSettingInfo `json:"setting"`
}

type SessionSettingInfo struct {
	IsSticky  bool `json:"is_sticky"`
	IsDeleted bool `json:"is_deleted"`
}

type CreateGroupRequest struct {
	GroupName string  `json:"group_name"`
	Members   []int64 `json:"members"`
}

type SessionService interface {
	OpenSession(openId, friendId int64) (SessionInfo, error)
	ListSession(userId int64) ([]SessionInfo, error)
	DeleteSession(userId int64, sessionId int64) error

	CreateGroupSession(userId int64, req CreateGroupRequest) (SessionInfo, error)
	AddSessionMember(sessionId int64, userId ...int64) error
	RemoveSessionMember(sessionId int64, userId int64) error
	ChangeMemberRole(ownerId int64, sessionId int64, userId int64, role int8) error
	DismissGroupSession(userId int64, sessionId int64) error

	UpdateSessionSetting(id int64, ss entity.SessionSetting) error
}

type sessionService struct {
	sessionEntity  entity.SessionEntity
	friendsService friends.FriendsService
}

func (s *sessionService) CreateGroupSession(userId int64, req CreateGroupRequest) (SessionInfo, error) {
	se := entity.Session{
		SessionName: req.GroupName,
		SessionType: entity.GroupSession,
		Members:     append(req.Members, userId),
	}
	tx := s.sessionEntity.Begin()
	id, err := s.sessionEntity.CreateSessionTx(tx, se, userId)
	if err != nil {
		tx.Rollback()
		return SessionInfo{}, err
	}
	err = s.sessionEntity.CreateSessionSettingTx(tx, entity.SessionSetting{SessionId: id, UserId: userId})
	if err != nil {
		tx.Rollback()
		return SessionInfo{}, err
	}
	tx.Commit()
	return SessionInfo{
		ID:          id,
		SessionName: req.GroupName,
		SessionType: entity.GroupSession,
		Members:     se.Members,
		Setting:     SessionSettingInfo{}}, nil
}

func (s *sessionService) OpenSession(openId, friendId int64) (SessionInfo, error) {
	if openId == friendId {
		return SessionInfo{}, errors.New(codec.OpenWithSelfErrCode)
	}
	ok, err := s.friendsService.IsFriend(openId, friendId)
	if err != nil {
		return SessionInfo{}, err
	}
	if !ok {
		return SessionInfo{}, errors.New(codec.NotFriendCode)
	}
	session, err := s.sessionEntity.GetSessionByUserId(openId, friendId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return SessionInfo{}, err
		}
		se := entity.Session{
			SessionName: strconv.FormatInt(friendId, 10),
			SessionType: entity.PointSession,
			Members:     []int64{openId, friendId},
		}
		tx := s.sessionEntity.Begin()
		id, err := s.sessionEntity.CreateSessionTx(tx, se, openId)
		if err != nil {
			tx.Rollback()
			return SessionInfo{}, err
		}
		err = s.sessionEntity.CreateSessionSettingTx(tx, entity.SessionSetting{SessionId: id, UserId: openId})
		if err != nil {
			tx.Rollback()
			return SessionInfo{}, err
		}
		tx.Commit()
		return SessionInfo{
			ID:          id,
			SessionName: se.SessionName,
			SessionType: se.SessionType,
			FriendId:    friendId,
			Setting:     SessionSettingInfo{}}, nil
	}
	return SessionInfo{
		ID:          session.ID,
		SessionName: session.SessionName,
		SessionType: session.SessionType,
		FriendId:    friendId,
		Setting:     SessionSettingInfo{}}, nil
}

func (s *sessionService) ListSession(userId int64) ([]SessionInfo, error) {
	session, err := s.sessionEntity.ListSession(userId)
	if err != nil {
		return nil, err
	}
	infoMap := make(map[int64]SessionInfo)
	for _, se := range session {
		settings, err := s.sessionEntity.GetSessionSetting(se.ID, userId)
		if err != nil {
			return nil, err
		}
		if settings.IsDeleted {
			continue
		}
		infoMap[se.ID] = SessionInfo{
			ID:          se.ID,
			SessionName: se.SessionName,
			SessionType: se.SessionType,
			Members:     make([]int64, 0),
			Setting:     SessionSettingInfo{IsSticky: settings.IsSticky, IsDeleted: settings.IsDeleted},
		}
	}
	if len(infoMap) == 0 {
		return nil, nil
	}

	sessionId := make([]int64, 0, len(session))
	for k := range infoMap {
		sessionId = append(sessionId, k)
	}
	sms, err := s.sessionEntity.ListSessionMember(sessionId...)
	if err != nil {
		return nil, err
	}
	for _, sm := range sms {
		info := infoMap[sm.ID]
		if sm.SessionType == entity.PointSession && sm.UserId != userId {
			info.FriendId = userId
		}
		if sm.SessionType == entity.GroupSession {
			info.Members = append(info.Members, sm.UserId)
		}
		infoMap[sm.SessionId] = info
	}
	infos := make([]SessionInfo, 0, len(infoMap))
	for _, v := range infoMap {
		infos = append(infos, v)
	}
	return infos, nil
}

func (s *sessionService) DeleteSession(userId int64, sessionId int64) error {
	sf, err := s.sessionEntity.GetSessionSetting(sessionId, userId)
	if err != nil {
		return err
	}
	sf.IsDeleted = true
	return s.UpdateSessionSetting(sf.ID, sf)
}

func (s *sessionService) DismissGroupSession(userId int64, sessionId int64) error {
	ok, err := s.sessionEntity.JudgeSessionRole(sessionId, userId, entity.RoleOwner)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New(codec.NotPermitDismissGroupCode)
	}
	return s.DeleteSession(userId, sessionId)
}

func (s *sessionService) AddSessionMember(sessionId int64, userId ...int64) error {
	return s.sessionEntity.CreateSessionMember(sessionId, userId...)
}

func (s *sessionService) RemoveSessionMember(sessionId int64, userId int64) error {
	return s.sessionEntity.RemoveSessionMember(sessionId, userId)
}

func (s *sessionService) ChangeMemberRole(ownerId int64, sessionId int64, userId int64, role int8) error {
	ok, err := s.sessionEntity.JudgeSessionRole(sessionId, ownerId, entity.RoleOwner)
	if err != nil {
		return err
	}
	if !ok {
		return errors.WithStackByCode(codec.NotPermitChangeRole)
	}
	return s.sessionEntity.ChangeMemberRole(sessionId, userId, role)
}

func (s *sessionService) UpdateSessionSetting(id int64, ss entity.SessionSetting) error {
	return s.sessionEntity.UpdateSessionSetting(id, ss)
}
