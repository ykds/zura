package services

import (
	"github.com/ykds/zura/internal/logic/entity"
	"github.com/ykds/zura/internal/logic/services/friend_application"
	"github.com/ykds/zura/internal/logic/services/friends"
	"github.com/ykds/zura/internal/logic/services/group"
	"github.com/ykds/zura/internal/logic/services/message"
	"github.com/ykds/zura/internal/logic/services/session"
	"github.com/ykds/zura/internal/logic/services/user"
	"github.com/ykds/zura/internal/logic/services/verify_code"
	"github.com/ykds/zura/pkg/cache"
	"github.com/ykds/zura/proto/comet"
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
	UserService              user.UserService
	FriendsService           friends.FriendsService
	FriendApplicationService friend_application.FriendApplicationService
	SessionService           session.SessionService
	VerifyCodeService        verify_code.VerifyCodeService
	MessageService           message.MessageService
	GroupService             group.GroupService
}

func NewServices(cache cache.Cache, entities *entity.Entity, cometClient comet.CometClient) {
	verifyCodeService := verify_code.NewVerifyCodeService(cache)
	friendService := friends.NewFriendsService(entities.FriendEntity, entities.UserEntity)
	sessionService := session.NewSessionService(entities.SessionEntity, entities.GroupEntity, entities.UserEntity)
	services = &Service{
		UserService:              user.NewUserService(cache, entities.UserEntity, verifyCodeService),
		FriendsService:           friendService,
		FriendApplicationService: friend_application.NewFriendApplicationService(cometClient, entities.FriendApplicationEntity, entities.FriendEntity),
		SessionService:           sessionService,
		VerifyCodeService:        verifyCodeService,
		MessageService:           message.NewMessageService(cometClient, entities.MessageEntity, entities.SessionEntity, entities.GroupEntity, entities.FriendEntity),
		GroupService:             group.NewGroupServer(entities.GroupEntity, entities.UserEntity, entities.SessionEntity),
	}
}
