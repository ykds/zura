package friend_applyment

import (
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
	UserId int64  `json:"user_id"`
	Markup string `json:"markup"`
}

const (
	ApplyTypeSend int8 = iota + 1
	ApplyTypeRecv
)

type Applyment struct {
	ID          int64  `json:"id"`
	UserId      int64  `json:"user_id"`
	Markup      string `json:"markup"`
	Type        int8   `json:"type"`
	Status      int8   `json:"status"`
	UpdatedTime int64  `json:"updated_time"`
}

type FriendApplymentService interface {
	ApplyFriend(userId int64, req ApplyRequest) error
	ListApplyments(userId int64) ([]Applyment, error)
	UpdateApplymentStatus(userId int64, id int64, status int8) error
	DeleteApplyment(id, userId int64) error
}

type friendApplymentService struct {
	friendApplymentEntity entity.FriendApplymentEntity
	friendEntity          entity.FriendEntity
}

func (f *friendApplymentService) ApplyFriend(userId int64, req ApplyRequest) error {
	// 如果已经是好友关系，返回”已是好友“错误。
	// 否则，查看是否有正在申请的记录，没有或者即使存在已过期或者拒绝的记录，都生成一条新的申请，有则更新记录并通知对方
	//
	// 除了 Apply 状态的只能存在一条记录，其他状态都可存在多条记录。
	ok, err := f.friendEntity.IsFriend(userId, req.UserId)
	if err != nil {
		return err
	}
	if ok {
		return errors.WithStackByCode(codec.HadBeFriendCode)
	}

	fa, err := f.friendApplymentEntity.GetApplyment(userId, req.UserId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return f.friendApplymentEntity.CreateApplyment(userId, req.UserId)
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
		var applyType int8
		if app.User1Id == userId {
			applyType = ApplyTypeSend
		} else {
			applyType = ApplyTypeRecv
		}
		a := Applyment{
			ID:          app.ID,
			UserId:      app.User2Id,
			Type:        applyType,
			Markup:      app.Markup,
			Status:      app.Status,
			UpdatedTime: app.UpdatedAt.Local().Unix(),
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

func (f *friendApplymentService) UpdateApplymentStatus(userId int64, id int64, status int8) error {
	if status != entity.Aggre && status != entity.Reject {
		return errors.WithStackByCode(codec.StatusErrCode)
	}
	fa, err := f.friendApplymentEntity.GetApplymentByID(id)
	if err != nil {
		return err
	}
	if fa.Status == entity.Aggre || fa.Status == entity.Reject {
		return errors.New(codec.DuplicateHandleApplymentErrCode)
	}
	if fa.User1Id == userId {
		return errors.WithStackByCode(codec.HandleSelfApplyErrCode)
	}
	// 添加好友，更新记录状态为通过
	if status == entity.Aggre {
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
