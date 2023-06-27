package entity

import (
	"github.com/ykds/zura/pkg/db"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/snowflake"
	"gorm.io/gorm"
)

const (
	RoleOwner int8 = iota + 1
	RoleManager
	RoleMember
)

type Group struct {
	BaseModel
	No      string `json:"no"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	OwnerId int64  `json:"owner_id"`
}

func (g Group) TableName() string {
	return "zura_group"
}

type GroupMember struct {
	BaseModel
	GroupId  int64  `json:"group_id"`
	UserId   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Role     int8   `json:"role"`
}

func (g GroupMember) TableName() string {
	return "zura_group_member"
}

func NewGroupEntity(db *db.Database) GroupEntity {
	return &groupEntity{
		baseEntity{db: db},
	}
}

type GroupEntity interface {
	Transaction
	GetGroup(id int64) (Group, error)
	ListGroups(userId int64) ([]Group, error)
	ListGroupById(id []int64) ([]Group, error)
	SearchGroup(where map[string]interface{}) ([]Group, error)
	CreateGroup(g *Group) (int64, error)
	CreateGroupTx(tx *gorm.DB, g *Group) (int64, error)
	UpdateGroup(id int64, g Group) error
	DeleteGroup(id int64) error

	GetGroupMember(groupId, memberId int64) (GroupMember, error)
	ListGroupMembers(groupId int64) ([]GroupMember, error)
	AddGroupMemberTx(tx *gorm.DB, member GroupMember) error
	AddGroupMember(member GroupMember) error
	RemoveGroupMember(groupId int64, memberId int64) error
	ChangeGroupMemberRole(groupId int64, memberId int64, role int8) error
	UpdateGroupMemberInfo(id int64, member GroupMember) error
	IsGroupMember(groupId int64, memberId int64) (bool, error)
}

type groupEntity struct {
	baseEntity
}

func (g2 groupEntity) SearchGroup(where map[string]interface{}) ([]Group, error) {
	g := make([]Group, 0)
	err := g2.db.Where(where).First(&g).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return g, err
}

func (g2 groupEntity) GetGroupMember(groupId, memberId int64) (GroupMember, error) {
	gm := GroupMember{}
	err := g2.db.First(&gm, "group_id=? AND user_id=?", groupId, memberId).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return gm, err
}

func (g2 groupEntity) IsGroupMember(groupId int64, memberId int64) (bool, error) {
	gm := GroupMember{}
	err := g2.db.Select("id").First(&gm, "group_id=? AND user_id=?", groupId, memberId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (g2 groupEntity) ListGroupById(id []int64) ([]Group, error) {
	groups := make([]Group, 0)
	err := g2.db.Find(&groups, "id IN ?", id).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return groups, err
}

func (g2 groupEntity) CreateGroupTx(tx *gorm.DB, g *Group) (int64, error) {
	err := tx.Transaction(func(t *gorm.DB) error {
		g.ID = snowflake.NewId()
		err := t.Create(g).Error
		if err != nil {
			return err
		}
		return t.Create(&GroupMember{
			GroupId: g.ID,
			UserId:  g.OwnerId,
			Role:    RoleOwner,
		}).Error
	})
	if err != nil {
		err = errors.WithStack(err)
	}
	return g.ID, err
}

func (g2 groupEntity) CreateGroup(g *Group) (int64, error) {
	return g2.CreateGroupTx(g2.db.DB, g)
}

func (g2 groupEntity) UpdateGroup(id int64, g Group) error {
	err := g2.db.Where("id = ?", id).Omit("owner_id").Updates(g).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (g2 groupEntity) DeleteGroup(id int64) error {
	err := g2.db.Transaction(func(tx *gorm.DB) error {
		err := g2.db.Delete(Group{}, "id=?", id).Error
		if err != nil {
			return err
		}
		return g2.db.Delete(GroupMember{}, "group_id=?", id).Error
	})
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (g2 groupEntity) ListGroups(userId int64) ([]Group, error) {
	gm := make([]GroupMember, 0)
	err := g2.db.Select("group_id").Find(&gm, "user_id=?", userId).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}
	groupIds := make([]int64, 0, len(gm))
	for _, item := range gm {
		groupIds = append(groupIds, item.GroupId)
	}
	result := make([]Group, 0)
	err = g2.db.Find(&result, "id IN ?", groupIds).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return result, err
}

func (g2 groupEntity) GetGroup(id int64) (Group, error) {
	g := Group{}
	err := g2.db.First(&g, "id=?", id).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return g, err
}

func (g2 groupEntity) AddGroupMemberTx(tx *gorm.DB, member GroupMember) error {
	_, err := g2.GetGroup(member.GroupId)
	if err != nil {
		return err
	}
	err = tx.Create(&member).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (g2 groupEntity) AddGroupMember(member GroupMember) error {
	return g2.AddGroupMemberTx(g2.db.DB, member)
}

func (g2 groupEntity) RemoveGroupMember(groupId int64, memberId int64) error {
	err := g2.db.Delete(GroupMember{}, "group_id=? AND user_id=?", groupId, memberId).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (g2 groupEntity) ChangeGroupMemberRole(groupId int64, memberId int64, role int8) error {
	err := g2.db.Model(GroupMember{}).Where("group_id=? AND user_id=?", groupId, memberId).Update("role", role).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (g2 groupEntity) UpdateGroupMemberInfo(id int64, member GroupMember) error {
	err := g2.db.Where("id=?", id).Select("nickname").Updates(&member).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (g2 groupEntity) ListGroupMembers(groupId int64) ([]GroupMember, error) {
	gm := make([]GroupMember, 0)
	err := g2.db.Find(&gm, "group_id=?", groupId).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return gm, err
}
