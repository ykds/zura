package friends

import "zura/internal/entity"

func NewFriendsService(friendsEntity entity.FriendEntity) FriendsService {
	return &friendsService{
		friendEntity: friendsEntity,
	}
}

type FriendsService interface {
}

type friendsService struct {
	friendEntity entity.FriendEntity
}