package entity

import (
	"fmt"
	"github.com/ykds/zura/pkg/db"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	c := db.DefaultConfig()
	c.Username = "root"
	c.Password = "zura123456"
	database := db.New(&c)

	err := database.Transaction(func(tx *gorm.DB) error {
		u := User{
			UserInfo: UserInfo{
				UpdatedEmailAt:    time.Now(),
				UpdatedPhoneAt:    time.Now(),
				UpdatedUsernameAt: time.Now(),
			},
		}
		err := tx.Create(&u).Error
		if err != nil {
			return err
		}
		fmt.Println(u.ID)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
