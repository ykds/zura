package entity

import (
	"zura/pkg/db"
	"zura/pkg/errors"

	"gorm.io/gorm"
)

const (
	Apply int8 = iota + 1
	Aggre
	Reject
	Expired
)

const (
	ApplymentNormal int8 = iota + 1
	ApplymentDeleteByOne
	ApplymentDeleteByTwo
)

type FriendApplyment struct {
	BaseModel
	User1Id   int64  `json:"user1_id"`
	User2Id   int64  `json:"user2_id"`
	Markup    string `json:"markup"`
	Status    int8   `json:"status"`
	DeletedBy int8   `json:"deleted_by"`
}

func (f FriendApplyment) TableName() string {
	return "zura_friend_applyment"
}

func NewFriendApplyment(db *db.Database) FriendApplymentEntity {
	return &friendApplymentEntity{
		baseEntity{db: db},
	}
}

type FriendApplymentEntity interface {
	Transaction
	GetApplyment(user1Id, user2Id int64) (FriendApplyment, error)
	GetApplymentByID(id int64) (FriendApplyment, error)
	CreateApplyment(user1Id, user2Id int64) error
	UpdateApplymentStatus(id int64, status int8) error
	UpdateApplymentStatusTx(tx *gorm.DB, id int64, status int8) error
	UpdateApplyment(id int64, fa FriendApplyment) error
	ListApplyments(userId int64) ([]FriendApplyment, error)
	DeleteApplyment(id int64, userId int64) error
}

type friendApplymentEntity struct {
	baseEntity
}

func (f *friendApplymentEntity) GetApplyment(user1Id int64, user2Id int64) (FriendApplyment, error) {
	fa := FriendApplyment{}
	err := f.db.First(&fa, "((user1_id=? AND user2_id=?) OR (user2_id=? AND user2_id=?)) AND status=?", user1Id, user2Id, user1Id, user2Id, Apply).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return fa, err
}

func (f *friendApplymentEntity) GetApplymentByID(id int64) (FriendApplyment, error) {
	fa := FriendApplyment{}
	err := f.db.First(&fa, "id=?", id).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return fa, err
}

func (f *friendApplymentEntity) CreateApplyment(user1Id, user2Id int64) error {
	err := f.db.Create(&FriendApplyment{User1Id: user1Id, User2Id: user2Id, Status: Apply, DeletedBy: ApplymentNormal}).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (f *friendApplymentEntity) UpdateApplymentStatus(id int64, status int8) error {
	return f.UpdateApplymentStatusTx(f.db.DB, id, status)
}

func (f *friendApplymentEntity) UpdateApplymentStatusTx(tx *gorm.DB, id int64, status int8) error {
	err := tx.Model(FriendApplyment{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (f *friendApplymentEntity) UpdateApplyment(id int64, fa FriendApplyment) error {
	err := f.db.Omit("user1_id", "user2_id", "status").Where("id=?", id).Updates(&fa).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (f *friendApplymentEntity) ListApplyments(userId int64) ([]FriendApplyment, error) {
	result := make([]FriendApplyment, 0)
	err := f.db.Raw("SELECT * from zura_friend_applyment WHERE user1_id=? AND deleted_by IN ? UNION ALL SELECT * FROM zura_friend_applyment WHERE user2_id=? AND deleted_by IN ?", userId, []int8{ApplymentNormal, ApplymentDeleteByTwo}, userId, []int8{ApplymentNormal, ApplymentDeleteByOne}).Order("updated_by desc").Scan(&result).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return result, err
}

func (f *friendApplymentEntity) DeleteApplyment(id int64, userId int64) error {
	fa := FriendApplyment{}
	err := f.db.First(&fa, "id=?", id).Error
	if err != nil {
		return errors.WithStack(err)
	}
	var deleteBy int8
	if fa.User1Id == userId {
		deleteBy = ApplymentDeleteByOne
		if fa.DeletedBy == ApplymentDeleteByTwo {
			err = f.db.Delete(&fa).Error
			if err != nil {
				return errors.WithStack(err)
			}
		}
	} else {
		deleteBy = ApplymentDeleteByTwo
		if fa.DeletedBy == ApplymentDeleteByOne {
			err = f.db.Delete(&fa).Error
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}
	fa.DeletedBy = deleteBy
	err = f.UpdateApplyment(fa.ID, fa)
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}
