package entity

import (
	"github.com/ykds/zura/internal/logic/codec"
	"github.com/ykds/zura/pkg/db"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/snowflake"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	PointSession int8 = iota + 1
	GroupSession
)

const (
	RoleOwner int8 = iota + 1
	RoleManager
	RoleMember
)

type Session struct {
	BaseModel
	SessionName string  `json:"sessionName"`
	SessionType int8    `json:"session_type"`
	Members     []int64 `json:"members" gorm:"-"`
}

func (r Session) TableName() string {
	return "zura_session"
}

type SessionSetting struct {
	BaseModel
	SessionId int64 `json:"session_id" gorm:"uniqueIndex:ids"`
	UserId    int64 `json:"user_id" gorm:"uniqueIndex:ids"`
	IsSticky  bool  `json:"is_sticky"`
	IsDeleted bool  `json:"is_deleted"`
}

func (r SessionSetting) TableName() string {
	return "zura_session_setting"
}

// 分表策略：
// 1. 创建 n 个相同结构的数据表
// 2. 修改程序，对新旧进行双写，旧表按原样写入，新表按sharding策略落到对应的表中；查询先查新表，如果没有再去旧表查询。
// 3. 启动一个异步程序，把第2步的修改上线后，把上线时间点之前的数据进行sharding到新表。
// 4. 当旧数据完成迁移后，去掉程序中查询旧表的逻辑。完成迁移。
type SessionMember struct {
	BaseModel
	UserId      int64 `json:"user_id" gorm:"uniqueIndex:ids"`
	SessionId   int64 `json:"session_id" gorm:"uniqueIndex:ids"`
	SessionType int8  `json:"session_type"`
	Role        int8  `json:"role"`
}

func (r SessionMember) TableName() string {
	return "zura_session_member"
}

func NewSessionEntity(db *db.Database) SessionEntity {
	return &sessionEntity{
		baseEntity{db: db},
	}
}

type SessionEntity interface {
	Transaction
	GetSession(id int64) (Session, error)
	GetSessionByUserId(user1Id, user2Id int64) (Session, error)
	CreateSession(session Session, creator int64) (int64, error)
	CreateSessionTx(tx *gorm.DB, session Session, creator int64) (int64, error)
	DeleteSession(id int64) error
	ListSession(userId int64) ([]Session, error)
	CreateSessionMember(sessionId int64, userId ...int64) error
	RemoveSessionMember(sessionId int64, userId int64) error
	ListSessionMember(sessionId ...int64) ([]SessionMember, error)
	ChangeMemberRole(sessionId int64, userId int64, role int8) error
	CreateSessionSetting(ss SessionSetting) error
	CreateSessionSettingTx(tx *gorm.DB, ss SessionSetting) error
	UpdateSessionSetting(id int64, ss SessionSetting) error
	DeleteSessionSetting(id int64) error
	GetSessionSetting(sessionId, userId int64) (SessionSetting, error)
	JudgeSessionRole(sessionId int64, userId int64, role int8) (bool, error)
}

type sessionEntity struct {
	baseEntity
}

func (r *sessionEntity) GetSession(id int64) (Session, error) {
	session := Session{}
	err := r.db.First(&session, "id=?", id).Error
	if err != nil {
		return session, errors.WithStack(err)
	}
	return session, nil
}

func (r *sessionEntity) GetSessionByUserId(user1Id, user2Id int64) (Session, error) {
	var result struct {
		SessionId int64
	}
	err := r.db.Raw("SELECT sm1.session_id FROM zura_session_member sm1 JOIN zura_session_member sm2 ON sm1.session_id = sm2.session_id WHERE sm1.user_id=? AND sm2.user_id=?", user1Id, user2Id).Scan(&result).Error
	if err != nil {
		return Session{}, errors.WithStack(err)
	}
	return r.GetSession(result.SessionId)
}

func (r *sessionEntity) CreateSession(session Session, createtor int64) (int64, error) {
	return r.CreateSessionTx(r.db.DB, session, createtor)
}

func (r *sessionEntity) CreateSessionTx(tx *gorm.DB, session Session, createtor int64) (int64, error) {
	err := tx.Transaction(func(t *gorm.DB) error {
		session.ID = snowflake.NewId()
		err := t.Create(&session).Error
		if err != nil {
			return err
		}
		members := make([]SessionMember, 0)
		for _, id := range session.Members {
			role := RoleMember
			if createtor == id && session.SessionType == GroupSession {
				role = RoleOwner
			}
			members = append(members, SessionMember{
				UserId:      id,
				SessionId:   session.ID,
				SessionType: session.SessionType,
				Role:        role,
			})
		}
		return t.Create(&members).Error
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return session.ID, nil
}

func (r *sessionEntity) DeleteSession(id int64) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Unscoped().Delete(Session{}, "id=?", id).Error
		if err != nil {
			return err
		}
		err = tx.Unscoped().Delete(SessionMember{}, "session_id=?", id).Error
		if err != nil {
			return err
		}
		return tx.Unscoped().Delete(SessionSetting{}, "session_id=?", id).Error
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *sessionEntity) ListSession(userId int64) ([]Session, error) {
	members := make([]SessionMember, 0)
	err := r.db.Select("session_id").Find(&members, "user_id=?", userId).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	rcList := make([]Session, 0)
	sessionId := make([]int64, 0)
	for _, member := range members {
		sessionId = append(sessionId, member.SessionId)
	}
	err = r.db.Find(&rcList, "id IN ?", sessionId).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return rcList, err
}

func (r *sessionEntity) ListSessionMember(sessionId ...int64) ([]SessionMember, error) {
	result := make([]SessionMember, 0)
	err := r.db.Find(&result, "session_id IN ?", sessionId).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return result, err
}

func (r *sessionEntity) CreateSessionMember(sessionId int64, userId ...int64) error {
	session, err := r.GetSession(sessionId)
	if err != nil {
		return err
	}
	if session.SessionType == PointSession {
		return errors.WithStackByCode(codec.AddMemberToPointSessionErrCode)
	}
	members := make([]SessionMember, 0)
	for _, id := range userId {
		members = append(members, SessionMember{SessionId: sessionId, UserId: id})
	}
	err = r.db.Clauses(clause.OnConflict{
		DoNothing: true,
	}).Create(&members).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *sessionEntity) RemoveSessionMember(sessionId int64, userId int64) error {
	session, err := r.GetSession(sessionId)
	if err != nil {
		return err
	}
	if session.SessionType == PointSession {
		return r.DeleteSession(sessionId)
	}
	err = r.db.Delete(&SessionMember{}, "session_id=? AND user_id=?", sessionId, userId).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *sessionEntity) ChangeMemberRole(sessionId int64, userId int64, role int8) error {
	var err error
	if role == RoleOwner {
		err = r.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("session_id = ? AND role=?", sessionId, RoleOwner).Update("role", RoleMember).Error; err != nil {
				return err
			}
			if err := tx.Where("session_id = ? AND user_id=?", sessionId, userId).Update("role", RoleOwner).Error; err != nil {
				return err
			}
			return nil
		})
	} else {
		err = r.db.Where("session_id = ? AND user_id=?", sessionId, userId).Update("role", role).Error
		if err != nil {
			err = errors.WithStack(err)
		}
	}
	return err
}

func (r *sessionEntity) CreateSessionSetting(ss SessionSetting) error {
	return r.CreateSessionSettingTx(r.db.DB, ss)
}

func (r *sessionEntity) CreateSessionSettingTx(tx *gorm.DB, ss SessionSetting) error {
	err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "session_id"}, {Name: "user_id"}},
		UpdateAll: true,
	}).Create(&ss).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *sessionEntity) UpdateSessionSetting(id int64, ss SessionSetting) error {
	err := r.db.Where("id=?", id).Omit("session_id", "user_id").Updates(&ss).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
func (r *sessionEntity) DeleteSessionSetting(id int64) error {
	err := r.db.Delete(&SessionSetting{}, "id=?", id).Error
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *sessionEntity) GetSessionSetting(sessionId, userId int64) (SessionSetting, error) {
	ss := SessionSetting{}
	err := r.db.First(&ss, "session_id=? AND user_id=?", sessionId, userId).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return ss, err
}

func (r *sessionEntity) JudgeSessionRole(sessionId int64, userId int64, role int8) (bool, error) {
	err := r.db.Where("session_id=? AND user_id=? AND role=?", sessionId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, errors.WithStack(err)
	}
	return true, nil
}
