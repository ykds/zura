package entity

import (
	"github.com/ykds/zura/pkg/db"
	"github.com/ykds/zura/pkg/errors"
	"gorm.io/gorm"
)

const (
	Apply int8 = iota + 1
	Agree
	Reject
	Expired
)

const (
	ApplicationNormal int8 = iota + 1
	ApplicationDeleteByOne
	ApplicationDeleteByTwo
)

type FriendApplication struct {
	BaseModel
	User1Id   int64  `json:"user1_id"`
	User2Id   int64  `json:"user2_id"`
	Markup    string `json:"markup"`
	Status    int8   `json:"status"`
	DeletedBy int8   `json:"deleted_by"`
}

func (f FriendApplication) TableName() string {
	return "zura_friend_application"
}

func NewFriendApplication(db *db.Database) FriendApplicationEntity {
	return &friendApplicationEntity{
		baseEntity{db: db},
	}
}

type FriendApplicationEntity interface {
	Transaction
	GetApplication(user1Id, user2Id int64) (FriendApplication, error)
	GetApplicationByID(id int64) (FriendApplication, error)
	CreateApplication(app FriendApplication) error
	UpdateApplicationStatus(id int64, status int8) error
	UpdateApplicationStatusTx(tx *gorm.DB, id int64, status int8) error
	UpdateApplication(id int64, fa FriendApplication) error
	ListApplications(userId int64) ([]FriendApplication, error)
	DeleteApplication(id int64, userId int64) error
}

type friendApplicationEntity struct {
	baseEntity
}

func (f *friendApplicationEntity) GetApplication(user1Id int64, user2Id int64) (FriendApplication, error) {
	fa := FriendApplication{}
	err := f.db.First(&fa, "((user1_id=? AND user2_id=?) OR (user2_id=? AND user2_id=?)) AND status=?", user1Id, user2Id, user1Id, user2Id, Apply).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return fa, err
}

func (f *friendApplicationEntity) GetApplicationByID(id int64) (FriendApplication, error) {
	fa := FriendApplication{}
	err := f.db.First(&fa, "id=?", id).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return fa, err
}

func (f *friendApplicationEntity) CreateApplication(app FriendApplication) error {
	app.Status = Apply
	app.DeletedBy = ApplicationNormal
	err := f.db.Create(&app).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (f *friendApplicationEntity) UpdateApplicationStatus(id int64, status int8) error {
	return f.UpdateApplicationStatusTx(f.db.DB, id, status)
}

func (f *friendApplicationEntity) UpdateApplicationStatusTx(tx *gorm.DB, id int64, status int8) error {
	err := tx.Model(FriendApplication{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (f *friendApplicationEntity) UpdateApplication(id int64, fa FriendApplication) error {
	err := f.db.Omit("user1_id", "user2_id", "status").Where("id=?", id).Updates(&fa).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}

func (f *friendApplicationEntity) ListApplications(userId int64) ([]FriendApplication, error) {
	result := make([]FriendApplication, 0)
	err := f.db.Raw("SELECT * from zura_friend_application WHERE user1_id=? AND deleted_by IN ? UNION ALL SELECT * FROM zura_friend_application WHERE user2_id=? AND deleted_by IN ?", userId, []int8{ApplicationNormal, ApplicationDeleteByTwo}, userId, []int8{ApplicationNormal, ApplicationDeleteByOne}).Order("updated_by desc").Scan(&result).Error
	if err != nil {
		err = errors.WithStack(err)
	}
	return result, err
}

func (f *friendApplicationEntity) DeleteApplication(id int64, userId int64) error {
	fa := FriendApplication{}
	err := f.db.First(&fa, "id=?", id).Error
	if err != nil {
		return errors.WithStack(err)
	}
	var deleteBy int8
	if fa.User1Id == userId {
		deleteBy = ApplicationDeleteByOne
		if fa.DeletedBy == ApplicationDeleteByTwo {
			err = f.db.Delete(&fa).Error
			if err != nil {
				return errors.WithStack(err)
			}
		}
	} else {
		deleteBy = ApplicationDeleteByTwo
		if fa.DeletedBy == ApplicationDeleteByOne {
			err = f.db.Delete(&fa).Error
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}
	fa.DeletedBy = deleteBy
	err = f.UpdateApplication(fa.ID, fa)
	if err != nil {
		err = errors.WithStack(err)
	}
	return err
}
