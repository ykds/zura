package entity

import (
	"github.com/ykds/zura/pkg/cache"
	"github.com/ykds/zura/pkg/db"
	"github.com/ykds/zura/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	PointSession int8 = iota + 1
	GroupSession
)

type UserSession struct {
	BaseModel
	SessionType int8  `json:"session_type"`
	SessionKey  int64 `json:"session_key"`
	UserId      int64 `json:"user_id" gorm:"index"`
	TargetId    int64 `json:"target_id" gorm:"index"`
	IsSticky    bool  `json:"is_sticky"`
}

func (r UserSession) TableName() string {
	return "zura_user_session"
}

func NewSessionEntity(cache cache.Cache, db *db.Database) SessionEntity {
	return &sessionEntity{
		cache:      cache,
		baseEntity: baseEntity{db: db},
	}
}

type SessionEntity interface {
	Transaction
	GetUserSession(where map[string]interface{}) (UserSession, error)
	GetUserSessionById(userId int64, sessionKey int64) (UserSession, error)
	ListSession(userId int64) ([]UserSession, error)
	CreateUserSessionTx(tx *gorm.DB, us UserSession) error
	UpdateUserSession(id int64, session UserSession) error
	DeleteUserSession(id int64) error
}

type sessionEntity struct {
	cache cache.Cache
	baseEntity
}

func (s2 sessionEntity) GetUserSession(where map[string]interface{}) (UserSession, error) {
	us := UserSession{}
	err := s2.db.Where(where).First(&us).Error
	if err != nil {
		return us, errors.WithStack(err)
	}
	return us, nil
}

func (s2 sessionEntity) ListSession(userId int64) ([]UserSession, error) {
	us := make([]UserSession, 0)
	err := s2.db.Find(&us, "user_id=?", userId).Order("is_sticky desc").Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return us, nil
}

func (s2 sessionEntity) CreateUserSessionTx(tx *gorm.DB, us UserSession) error {
	err := tx.Unscoped().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "session_type"}, {Name: "user_id"}, {Name: "target_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"deleted_at": nil}),
	}).Create(&us).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (s2 sessionEntity) UpdateUserSession(id int64, session UserSession) error {
	err := s2.db.Where("id=?", id).Omit("session_type", "user_id", "target_id", "session_key").Updates(&session).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (s2 sessionEntity) DeleteUserSession(id int64) error {
	err := s2.db.Delete(UserSession{}, "id=?", id).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (s2 sessionEntity) GetUserSessionById(userId, sessionKey int64) (UserSession, error) {
	us := UserSession{}
	err := s2.db.Where("user_id = ? AND session_key=?", userId, sessionKey).First(&us).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return us, err
}
