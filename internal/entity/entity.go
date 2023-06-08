package entity

import (
	"time"
	"zura/pkg/cache"
	"zura/pkg/db"

	"gorm.io/gorm"
)

var (
	entity *Entity
)

type BaseModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

var tables = []interface{}{
	User{}, Friends{}, RecentContacts{},
}

func migrateTable(d *db.Database) error {
	return d.AutoMigrate(tables...)
}

type Entity struct {
	UserEntity           UserEntity
	FriendEntity         FriendEntity
	RecentContactsEntity RecentContactsEntity
}

func GetEntity() *Entity {
	if entity == nil {
		panic("never init entity")
	}
	return entity
}

func NewEntity(database *db.Database, cache *cache.Redis) {
	if err := migrateTable(database); err != nil {
		panic(err)
	}
	entity = &Entity{
		UserEntity:           NewUserEntity(database),
		FriendEntity:         NewFriendEntity(database),
		RecentContactsEntity: NewRecentContactsEntity(database),
	}
}
