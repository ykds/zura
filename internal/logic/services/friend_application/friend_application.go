package friend_application

import (
	"github.com/ykds/zura/internal/logic/codec"
	"github.com/ykds/zura/internal/logic/entity"
	"github.com/ykds/zura/pkg/errors"
	"gorm.io/gorm"
)

func NewFriendApplicationService(friendApplicationEntity entity.FriendApplicationEntity, friendEntity entity.FriendEntity) FriendApplicationService {
	return &friendApplicationService{
		friendApplicationEntity: friendApplicationEntity,
		friendEntity:            friendEntity,
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

type Application struct {
	ID          int64  `json:"id"`
	UserId      int64  `json:"user_id"`
	Markup      string `json:"markup"`
	Type        int8   `json:"type"`
	Status      int8   `json:"status"`
	UpdatedTime int64  `json:"updated_time"`
}

type FriendApplicationService interface {
	ApplyFriend(userId int64, req ApplyRequest) error
	ListApplications(userId int64) ([]Application, error)
	UpdateApplicationStatus(userId int64, id int64, status int8) error
	DeleteApplication(id, userId int64) error
}

type friendApplicationService struct {
	friendApplicationEntity entity.FriendApplicationEntity
	friendEntity            entity.FriendEntity
}

func (f *friendApplicationService) ApplyFriend(userId int64, req ApplyRequest) error {
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

	fa, err := f.friendApplicationEntity.GetApplication(userId, req.UserId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return f.friendApplicationEntity.CreateApplication(userId, req.UserId)
		}
		err = nil
	}
	err = f.friendApplicationEntity.UpdateApplication(fa.ID, entity.FriendApplication{Markup: req.Markup, DeletedBy: entity.ApplicationNormal})
	if err != nil {
		return err
	}
	// TODO 通知对方
	return nil
}

func (f *friendApplicationService) ListApplications(userId int64) ([]Application, error) {
	apps, err := f.friendApplicationEntity.ListApplications(userId)
	if err != nil {
		return nil, err
	}
	result := make([]Application, 0)
	for _, app := range apps {
		var applyType int8
		if app.User1Id == userId {
			applyType = ApplyTypeSend
		} else {
			applyType = ApplyTypeRecv
		}
		a := Application{
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

func (f *friendApplicationService) UpdateApplicationStatus(userId int64, id int64, status int8) error {
	if status != entity.Aggre && status != entity.Reject {
		return errors.WithStackByCode(codec.StatusErrCode)
	}
	fa, err := f.friendApplicationEntity.GetApplicationByID(id)
	if err != nil {
		return err
	}
	if fa.Status == entity.Aggre || fa.Status == entity.Reject {
		return errors.New(codec.DuplicateHandleApplicationErrCode)
	}
	if fa.User1Id == userId {
		return errors.WithStackByCode(codec.HandleSelfApplyErrCode)
	}
	// 添加好友，更新记录状态为通过
	if status == entity.Aggre {
		tx := f.friendApplicationEntity.Begin()
		err = f.friendEntity.AddFriendTx(tx, fa.User1Id, fa.User2Id)
		if err != nil {
			f.friendApplicationEntity.Rollback(tx)
			return err
		}
		err = f.friendApplicationEntity.UpdateApplicationStatusTx(tx, id, status)
		if err != nil {
			f.friendApplicationEntity.Rollback(tx)
			return err
		}
		f.friendApplicationEntity.Commit(tx)
	}
	return f.friendApplicationEntity.UpdateApplicationStatus(id, status)
}

func (f *friendApplicationService) DeleteApplication(id int64, userId int64) error {
	return f.friendApplicationEntity.DeleteApplication(id, userId)
}
