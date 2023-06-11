package services

import (
	"zura/internal/logic/entity"
	"zura/internal/logic/services/friend_applyment"
	"zura/internal/logic/services/friends"
	"zura/internal/logic/services/session"
	"zura/internal/logic/services/user"
	"zura/internal/logic/services/verify_code"
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
	UserService            user.UserService
	FriendsService         friends.FriendsService
	FriendApplymentService friend_applyment.FriendApplymentService
	SessionService         session.SessionService
	VerifyCodeService      verify_code.VerifyCodeService
}

func NewServices(cache *cache.Redis, entities *entity.Entity) {
	verifyCodeService := verify_code.NewVerifyCodeService(cache)
	friendService := friends.NewFriendsService(entities.FriendEntity, entities.UserEntity)
	services = &Service{
		UserService:            user.NewUserService(entities.UserEntity, verifyCodeService),
		FriendsService:         friendService,
		FriendApplymentService: friend_applyment.NewFriendApplymentService(entities.FriendApplymentEntity, entities.FriendEntity),
		SessionService:         session.NewSessionService(entities.SessionEntity, friendService),
		VerifyCodeService:      verify_code.NewVerifyCodeService(cache),
	}
}
