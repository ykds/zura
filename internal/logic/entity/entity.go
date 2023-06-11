package entity

import (
	"time"
	"zura/pkg/cache"
	"zura/pkg/db"
	"zura/pkg/snowflake"

	"gorm.io/gorm"
)

var (
	entity *Entity
)

type BaseModel struct {
	ID        int64          `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == 0 {
		b.ID = snowflake.NewId()
	}
	return
}

var tables = []interface{}{
	User{}, Friends{}, Session{}, FriendApplyment{}, SessionSetting{}, SessionMember{},
}

func migrateTable(d *db.Database) error {
	return d.AutoMigrate(tables...)
}

type Entity struct {
	UserEntity            UserEntity
	FriendEntity          FriendEntity
	FriendApplymentEntity FriendApplymentEntity
	SessionEntity         SessionEntity
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
		UserEntity:            NewUserEntity(database),
		FriendEntity:          NewFriendEntity(database),
		FriendApplymentEntity: NewFriendApplyment(database),
		SessionEntity:         NewSessionEntity(database),
	}
}

var _ Transaction = &baseEntity{}

type Transaction interface {
	Begin() *gorm.DB
	Commit(*gorm.DB) *gorm.DB
	Rollback(*gorm.DB) *gorm.DB
}

type baseEntity struct {
	db *db.Database
}

func (b *baseEntity) Begin() *gorm.DB {
	return b.db.Begin()
}

func (b *baseEntity) Commit(tx *gorm.DB) *gorm.DB {
	return tx.Commit()
}

func (b *baseEntity) Rollback(tx *gorm.DB) *gorm.DB {
	return tx.Rollback()
}
