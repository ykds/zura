package entity

import (
	"sync"
	"zira/internal/config"
	"zira/internal/entity/user"
	"zira/pkg/db"
)

var (
	database *db.Database
	once sync.Once
	entity *Entity
)

var Tables = []interface{}{
	user.User{}, user.Friends{}, user.RecentContacts{},
}

func GetDB() *db.Database {
	once.Do(func() {
		database = db.New(db.WithConfig(&config.GetConfig().Database))
	})
	return database
}

func MigrateTable() error {
	return database.AutoMigrate(Tables...)
}

type Entity struct {
}

func NewEntity() error {
	entity = &Entity{}
	return nil
}

func GetEntity() *Entity {
	return entity
}