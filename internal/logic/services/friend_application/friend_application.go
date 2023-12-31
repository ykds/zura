package friend_application

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/internal/logic/codec"
	"github.com/ykds/zura/internal/logic/entity"
	"github.com/ykds/zura/pkg/cache"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/kafka"
	"github.com/ykds/zura/pkg/log"
	"github.com/ykds/zura/proto/logic"
	"github.com/ykds/zura/proto/protocol"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func NewFriendApplicationService(cache cache.Cache, kafkaProducer *kafka.Producer, friendApplicationEntity entity.FriendApplicationEntity, friendEntity entity.FriendEntity) FriendApplicationService {
	return &friendApplicationService{
		cache:                   cache,
		kafkaProducer:           kafkaProducer,
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
	ListNewApplications(userId int64) ([]Application, error)
	UpdateApplicationStatus(userId int64, id int64, status int8) error
	DeleteApplication(id, userId int64) error
}

type friendApplicationService struct {
	cache                   cache.Cache
	kafkaProducer           *kafka.Producer
	friendApplicationEntity entity.FriendApplicationEntity
	friendEntity            entity.FriendEntity
}

func (f *friendApplicationService) ListNewApplications(userId int64) ([]Application, error) {
	applications, err := f.ListApplications(userId)
	if err != nil {
		return nil, err
	}
	newApp := make([]Application, 0, len(applications))
	for _, app := range applications {
		if app.Status != entity.Apply {
			continue
		}
		newApp = append(newApp, app)
	}
	return newApp, nil
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
			return err
		}
		err = f.friendApplicationEntity.CreateApplication(entity.FriendApplication{
			User1Id: userId,
			User2Id: req.UserId,
			Markup:  req.Markup,
		})
		if err != nil {
			return err
		}
	} else {
		err = f.friendApplicationEntity.UpdateApplication(fa.ID, entity.FriendApplication{Markup: req.Markup, DeletedBy: entity.ApplicationNormal})
		if err != nil {
			return err
		}
	}
	result, err := f.cache.Get(context.Background(), fmt.Sprintf(common.UserOnlineCacheKey, userId))
	if err != nil {
		return err
	}
	serverId, _ := strconv.ParseInt(result, 10, 64)
	msg := &logic.PushNotification{Op: protocol.OpNewApplication, ToUserId: []int64{userId}, Server: int32(serverId)}
	marshal, _ := json.Marshal(msg)
	err = f.kafkaProducer.WriteMessage(context.TODO(), "", marshal)
	if err != nil {
		log.Errorf("push notification error: %v", err)
	}
	return err
}

func (f *friendApplicationService) ListApplications(userId int64) ([]Application, error) {
	apps, err := f.friendApplicationEntity.ListApplications(userId)
	if err != nil {
		return nil, err
	}
	expiredApp := make([]entity.FriendApplication, 0)
	result := make([]Application, 0)
	for _, app := range apps {
		// 过期
		if app.UpdatedAt.Local().Add(24 * time.Hour * 3).After(time.Now()) {
			expiredApp = append(expiredApp, app)
			continue
		}
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
	go func() {
		for _, app := range expiredApp {
			_ = f.friendApplicationEntity.UpdateApplicationStatus(app.ID, entity.Expired)
		}
	}()
	return result, nil
}

func (f *friendApplicationService) UpdateApplicationStatus(userId int64, id int64, status int8) error {
	if status != entity.Agree && status != entity.Reject {
		return errors.WithStackByCode(codec.StatusErrCode)
	}
	fa, err := f.friendApplicationEntity.GetApplicationByID(id)
	if err != nil {
		return err
	}
	if fa.Status == entity.Expired {
		return errors.New(codec.ExpiredCode)
	}
	if fa.Status == entity.Agree || fa.Status == entity.Reject {
		return errors.New(codec.DuplicateHandleApplicationErrCode)
	}
	if fa.User1Id == userId {
		return errors.WithStackByCode(codec.HandleSelfApplyErrCode)
	}
	// 添加好友，更新记录状态为通过
	tx := f.friendApplicationEntity.Begin()
	if status == entity.Agree {
		err = f.friendEntity.AddFriendTx(tx, fa.User1Id, fa.User2Id)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = f.friendApplicationEntity.UpdateApplicationStatusTx(tx, id, status)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		err = f.friendApplicationEntity.UpdateApplicationStatusTx(tx, id, status)
		if err != nil {
			return err
		}
	}
	result, err := f.cache.Get(context.Background(), fmt.Sprintf(common.UserOnlineCacheKey, userId))
	if err != nil {
		tx.Rollback()
		return err
	}
	serverId, _ := strconv.ParseInt(result, 10, 64)
	body, _ := json.Marshal(map[string]interface{}{"id": id, "status": status})
	msg := &logic.PushNotification{Op: protocol.OpApplicationHandlerResult, ToUserId: []int64{userId}, Server: int32(serverId), Body: body}
	marshal, _ := json.Marshal(msg)
	err = f.kafkaProducer.WriteMessage(context.TODO(), "", marshal)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (f *friendApplicationService) DeleteApplication(id int64, userId int64) error {
	return f.friendApplicationEntity.DeleteApplication(id, userId)
}
