package services

import (
	"zura/internal/entity"
	"zura/internal/services/friends"
	"zura/internal/services/recent_contacts"
	"zura/internal/services/user"
	"zura/internal/services/verify_code"
	"zura/pkg/cache"
)

var (
	services *Service
)

func GetServices() *Service {
	if services == nil {
		panic("未初始化Service")
	}
	return services
}

type Service struct {
	UserService           user.UserService
	FriendsService        friends.FriendsService
	RecentContactsService recent_contacts.RecentContactsService
	VerifyCodeService     verify_code.VerifyCodeService
}

func NewServices(cache *cache.Redis, entities *entity.Entity) {
	verifyCodeService := verify_code.NewVerifyCodeService(cache)
	services = &Service{
		UserService:           user.NewUserService(entities.UserEntity, verifyCodeService),
		FriendsService:        friends.NewFriendsService(entities.FriendEntity),
		RecentContactsService: recent_contacts.NewRecentContactsService(entities.RecentContactsEntity),
	}
}
