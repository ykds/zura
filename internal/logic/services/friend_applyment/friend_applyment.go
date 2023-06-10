package friend_applyment

import (
	"time"
	"zura/internal/logic/codec"
	"zura/internal/logic/entity"
	"zura/pkg/errors"

	"gorm.io/gorm"
)

func NewFriendApplymentService(friendApplymentEntity entity.FriendApplymentEntity, friendEntity entity.FriendEntity) FriendApplymentService {
	return &friendApplymentService{
		friendApplymentEntity: friendApplymentEntity,
		friendEntity:          friendEntity,
	}
}

type ApplyRequest struct {
	User1Id int64  `json:"user1_id"`
	User2Id int64  `json:"user2_id"`
	Markup  string `json:"markup"`
}

type Applyment struct {
	UserId    int64     `json:"user_id"`
	Markup    string    `json:"markup"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FriendApplymentService interface {
	ApplyFriend(req ApplyRequest) error
	ListApplyments(userId int64) ([]Applyment, error)
	UpdateApplymentStatus(id int64, status int8) error
	DeleteApplyment(id, userId int64) error
}

type friendApplymentService struct {
	friendApplymentEntity entity.FriendApplymentEntity
	friendEntity          entity.FriendEntity
}

func (f *friendApplymentService) ApplyFriend(req ApplyRequest) error {
	// 如果已经是好友关系，返回”已是好友“错误。
	// 否则，查看是否有正在申请的记录，没有或者即使存在已过期或者拒绝的记录，都生成一条新的申请，有则更新记录并通知对方
	//
	// 除了 Apply 状态的只能存在一条记录，其他状态都可存在多条记录。
	ok, err := f.friendEntity.IsFriend(req.User1Id, req.User2Id)
	if err != nil {
		return err
	}
	if ok {
		return errors.WithStackByCode(codec.HadBeFriendCode)
	}

	fa, err := f.friendApplymentEntity.GetApplyment(req.User1Id, req.User2Id)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return f.friendApplymentEntity.CreateApplyment(req.User1Id, req.User2Id)
		}
		err = nil
	}
	err = f.friendApplymentEntity.UpdateApplyment(fa.ID, entity.FriendApplyment{Markup: req.Markup, DeletedBy: entity.ApplymentNormal})
	if err != nil {
		return err
	}
	// TODO 通知对方
	return nil
}

func (f *friendApplymentService) ListApplyments(userId int64) ([]Applyment, error) {
	apps, err := f.friendApplymentEntity.ListApplyments(userId)
	if err != nil {
		return nil, err
	}
	result := make([]Applyment, 0)
	for _, app := range apps {
		a := Applyment{
			UserId:    app.User2Id,
			Markup:    app.Markup,
			UpdatedAt: app.UpdatedAt,
		}
		if app.User1Id == userId {
			a.UserId = app.User2Id
		} else {
			a.UserId = app.User1Id
		}
		result = append(result, a)
	}
	return result, nil
}

func (f *friendApplymentService) UpdateApplymentStatus(id int64, status int8) error {
	if status != entity.Aggre && status != entity.Reject {
		return errors.WithStackByCode(codec.StatusErrCode)
	}
	// 添加好友，更新记录状态为通过
	if status == entity.Aggre {
		fa, err := f.friendApplymentEntity.GetApplymentByID(id)
		if err != nil {
			return err
		}
		tx := f.friendApplymentEntity.Begin()
		err = f.friendEntity.AddFriendTx(tx, fa.User1Id, fa.User2Id)
		if err != nil {
			f.friendApplymentEntity.Rollback(tx)
			return err
		}
		err = f.friendApplymentEntity.UpdateApplymentStatusTx(tx, id, status)
		if err != nil {
			f.friendApplymentEntity.Rollback(tx)
			return err
		}
		f.friendApplymentEntity.Commit(tx)
	}
	return f.friendApplymentEntity.UpdateApplymentStatus(id, status)
}

func (f *friendApplymentService) DeleteApplyment(id int64, userId int64) error {
	return f.friendApplymentEntity.DeleteApplyment(id, userId)
}
