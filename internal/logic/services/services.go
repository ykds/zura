package services

import (
	kafka2 "github.com/segmentio/kafka-go"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/internal/logic/balancer"
	"github.com/ykds/zura/internal/logic/entity"
	"github.com/ykds/zura/internal/logic/services/friend_application"
	"github.com/ykds/zura/internal/logic/services/friends"
	"github.com/ykds/zura/internal/logic/services/group"
	"github.com/ykds/zura/internal/logic/services/message"
	"github.com/ykds/zura/internal/logic/services/session"
	"github.com/ykds/zura/internal/logic/services/user"
	"github.com/ykds/zura/internal/logic/services/verify_code"
	"github.com/ykds/zura/pkg/cache"
	"github.com/ykds/zura/pkg/kafka"
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
	UserService              user.Service
	FriendsService           friends.Service
	FriendApplicationService friend_application.FriendApplicationService
	SessionService           session.Service
	VerifyCodeService        verify_code.VerifyCodeService
	MessageService           message.Service
	GroupService             group.Service
}

func NewServices(cache cache.Cache, entities *entity.Entity, kafkaManager *kafka.Kafka) {
	verifyCodeService := verify_code.NewVerifyCodeService(cache)
	friendService := friends.NewFriendsService(entities.FriendEntity, entities.UserEntity)
	sessionService := session.NewSessionService(entities.SessionEntity, entities.GroupEntity, entities.UserEntity)

	producer := kafkaManager.NewProducer(common.MessageTopic, kafka.WithCustomBalancer(kafka2.BalancerFunc(balancer.SessionIdBalance)))
	notifyProducer := kafkaManager.NewProducer(common.NotificationTopic)
	services = &Service{
		UserService:              user.NewUserService(cache, entities.UserEntity, verifyCodeService),
		FriendsService:           friendService,
		FriendApplicationService: friend_application.NewFriendApplicationService(cache, notifyProducer, entities.FriendApplicationEntity, entities.FriendEntity),
		SessionService:           sessionService,
		VerifyCodeService:        verifyCodeService,
		MessageService:           message.NewMessageService(cache, producer, entities.MessageEntity, entities.SessionEntity, entities.GroupEntity, entities.FriendEntity),
		GroupService:             group.NewGroupServer(entities.GroupEntity, entities.UserEntity, entities.SessionEntity),
	}
}
