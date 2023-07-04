package group

import (
	"github.com/ykds/zura/internal/logic/codec"
	"github.com/ykds/zura/internal/logic/entity"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/random"
)

func NewGroupServer(groupEntity entity.GroupEntity, userEntity entity.UserEntity, sessionEntity entity.SessionEntity) Service {
	return &groupService{
		groupEntity:   groupEntity,
		userEntity:    userEntity,
		sessionEntity: sessionEntity,
	}
}

type CreateGroupRequest struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type Info struct {
	ID      int64  `json:"id"`
	No      string `json:"no"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	OwnerId int64  `json:"owner_id"`
}

type UpdateGroupRequest struct {
	GroupId int64  `json:"group_id"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
}

type AddGroupMemberRequest struct {
	GroupId int64 `json:"group_id"`
	UserId  int64 `json:"user_id"`
	Role    int8  `json:"role"`
}

type RemoveGroupMemberRequest struct {
	GroupId int64 `json:"group_id"`
	UserId  int64 `json:"user_id"`
}

type UpdateMemberInfoRequest struct {
	GroupId  int64  `json:"group_id"`
	NickName string `json:"nick_name"`
}

type MemberInfo struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	NickName string `json:"nick_name"`
	Role     int8   `json:"role"`
}

type ChangeMemberRoleRequest struct {
	GroupId int64 `json:"group_id"`
	UserId  int64 `json:"user_id"`
	Role    int8  `json:"role"`
}

type SearchGroupRequest struct {
	No   string `form:"no"`
	Name string `form:"name"`
}

type Service interface {
	CreateGroup(userId int64, req CreateGroupRequest) error
	ListGroup(userId int64) ([]Info, error)
	UpdateGroup(userId int64, req UpdateGroupRequest) error
	DismissGroup(userId int64, groupId int64) error

	AddGroupMember(userId int64, req AddGroupMemberRequest) error
	RemoveGroupMember(userId int64, req RemoveGroupMemberRequest) error
	UpdateMemberInfo(userId int64, req UpdateMemberInfoRequest) error
	ListGroupMembers(userId int64, groupId int64) ([]MemberInfo, error)
	ChangeMemberRole(userId int64, req ChangeMemberRoleRequest) error
	SearchGroup(req SearchGroupRequest) ([]Info, error)
}

type groupService struct {
	groupEntity   entity.GroupEntity
	userEntity    entity.UserEntity
	sessionEntity entity.SessionEntity
}

func (g groupService) SearchGroup(req SearchGroupRequest) ([]Info, error) {
	where := make(map[string]interface{})
	if req.No != "" {
		where["no"] = req.No
	}
	if req.Name != "" {
		where["name"] = "%" + req.Name + "%"
	}
	group, err := g.groupEntity.SearchGroup(where)
	if err != nil {
		return nil, err
	}
	infos := make([]Info, 0, len(group))
	for _, item := range group {
		infos = append(infos, Info{
			ID:      item.ID,
			No:      item.No,
			Name:    item.Name,
			Avatar:  item.Avatar,
			OwnerId: item.OwnerId,
		})
	}
	return infos, nil
}

func (g groupService) CreateGroup(userId int64, req CreateGroupRequest) error {
	group := entity.Group{
		Name:    req.Name,
		Avatar:  req.Avatar,
		OwnerId: userId,
		No:      random.RandNum(8),
	}
	tx := g.groupEntity.Begin()
	_, err := g.groupEntity.CreateGroupTx(tx, &group)
	if err != nil {
		return err
	}
	err = g.sessionEntity.CreateUserSessionTx(tx, entity.UserSession{
		SessionType: entity.GroupSession,
		UserId:      userId,
		TargetId:    group.ID,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func (g groupService) ListGroup(userId int64) ([]Info, error) {
	group, err := g.ListGroup(userId)
	if err != nil {
		return nil, err
	}
	result := make([]Info, 0, len(group))
	for _, g := range group {
		result = append(result, Info{
			ID:      g.ID,
			No:      g.No,
			Name:    g.Name,
			Avatar:  g.Avatar,
			OwnerId: g.OwnerId,
		})
	}
	return result, nil
}

func (g groupService) UpdateGroup(userId int64, req UpdateGroupRequest) error {
	group, err := g.groupEntity.GetGroup(req.GroupId)
	if err != nil {
		return err
	}
	if group.OwnerId != userId {
		return errors.WithStackByCode(codec.NotPermitCode)
	}
	return g.groupEntity.UpdateGroup(req.GroupId, entity.Group{
		Name:   req.Name,
		Avatar: req.Avatar,
	})
}

func (g groupService) DismissGroup(userId int64, groupId int64) error {
	group, err := g.groupEntity.GetGroup(groupId)
	if err != nil {
		return err
	}
	if group.OwnerId != userId {
		return errors.WithStackByCode(codec.NotPermitCode)
	}
	return g.groupEntity.DeleteGroup(groupId)
}

func (g groupService) AddGroupMember(userId int64, req AddGroupMemberRequest) error {
	// 判断添加成员的操作者是否有权限添加
	member, err := g.groupEntity.GetGroupMember(req.GroupId, userId)
	if err != nil {
		return err
	}
	if member.Role == entity.RoleMember {
		return errors.WithStackByCode(codec.NotPermitCode)
	}
	// 判断要添加的人是否已经是成员
	ok, err := g.groupEntity.IsGroupMember(req.GroupId, req.UserId)
	if err != nil {
		return err
	}
	if ok {
		return errors.WithStackByCode(codec.HadAddGroupCode)
	}
	if req.Role != entity.RoleManager && req.Role != entity.RoleMember {
		return errors.WithStackByCode(codec.UnSupportRoleCode)
	}
	tx := g.groupEntity.Begin()
	err = g.groupEntity.AddGroupMemberTx(tx, entity.GroupMember{
		GroupId: req.GroupId,
		UserId:  req.UserId,
		Role:    req.Role,
	})
	if err != nil {
		return err
	}
	err = g.sessionEntity.CreateUserSessionTx(tx, entity.UserSession{
		SessionType: entity.GroupSession,
		UserId:      req.UserId,
		TargetId:    req.GroupId,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (g groupService) RemoveGroupMember(userId int64, req RemoveGroupMemberRequest) error {
	member, err := g.groupEntity.GetGroupMember(req.GroupId, userId)
	if err != nil {
		return err
	}
	if member.Role == entity.RoleMember {
		return errors.WithStackByCode(codec.NotPermitCode)
	}
	return g.groupEntity.RemoveGroupMember(req.GroupId, req.UserId)
}

func (g groupService) UpdateMemberInfo(userId int64, req UpdateMemberInfoRequest) error {
	return g.groupEntity.UpdateGroupMemberInfo(userId, entity.GroupMember{
		GroupId:  req.GroupId,
		Nickname: req.NickName,
	})
}

func (g groupService) ListGroupMembers(userId int64, groupId int64) ([]MemberInfo, error) {
	ok, err := g.groupEntity.IsGroupMember(groupId, userId)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.WithStackByCode(codec.NotGroupMember)
	}
	members, err := g.groupEntity.ListGroupMembers(groupId)
	if err != nil {
		return nil, err
	}
	result := make([]MemberInfo, 0, len(members))
	for _, item := range members {
		user, err := g.userEntity.GetUserById(item.UserId)
		if err != nil {
			return nil, err
		}
		result = append(result, MemberInfo{
			UserId:   item.UserId,
			Username: user.Username,
			Avatar:   user.Avatar,
			NickName: item.Nickname,
			Role:     item.Role,
		})
	}
	return result, nil
}

func (g groupService) ChangeMemberRole(userId int64, req ChangeMemberRoleRequest) error {
	member, err := g.groupEntity.GetGroupMember(req.GroupId, userId)
	if err != nil {
		return err
	}
	if member.Role != entity.RoleOwner {
		return errors.WithStackByCode(codec.NotPermitCode)
	}
	if req.Role != entity.RoleManager && req.Role != entity.RoleMember {
		return errors.WithStackByCode(codec.UnSupportRoleCode)
	}
	return g.groupEntity.ChangeGroupMemberRole(req.GroupId, req.UserId, req.Role)
}
