package friends

import (
	"github.com/ykds/zura/internal/logic/codec"
	"github.com/ykds/zura/internal/logic/entity"
	"github.com/ykds/zura/pkg/errors"
)

func NewFriendsService(friendsEntity entity.FriendEntity, userEntity entity.UserEntity) Service {
	return &friendsService{
		friendEntity: friendsEntity,
		userEntity:   userEntity,
	}
}

type FriendInfo struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type Service interface {
	AddFriend(user1Id, user2Id int64) error
	DeleteFriend(user1Id, user2Id int64) error
	ListFriends(userId int64) ([]FriendInfo, error)
	IsFriend(user1Id, user2Id int64) (bool, error)
}

type friendsService struct {
	friendEntity entity.FriendEntity
	userEntity   entity.UserEntity
}

func (f *friendsService) AddFriend(user1Id int64, user2Id int64) error {
	friend, err := f.friendEntity.GetFriend(user1Id, user2Id)
	if err != nil {
		return err
	}
	if friend.Status != entity.Normal {
		return errors.WithStackByCode(codec.HadBeFriendCode)
	}
	return f.friendEntity.AddFriend(user1Id, user2Id)
}

func (f *friendsService) DeleteFriend(user1Id int64, user2Id int64) error {
	friend, err := f.friendEntity.GetFriend(user1Id, user2Id)
	if err != nil {
		return err
	}
	// 已互删，忽略
	if friend.Status == entity.DeletedByEach {
		return nil
	}
	// 仅当前状态为，正常，我删了对方，对方删了我 的情况下才可以进行删除好友操作
	if friend.Status != entity.Normal && friend.Status != entity.DeletedByOne && friend.Status != entity.DeletedByTwo {
		return nil
	}
	// 如何对方已经删了我，应该把状态设置为互删
	var status int8
	if friend.User1Id == user1Id {
		status = entity.DeletedByOne
		if friend.Status == entity.DeletedByTwo {
			status = entity.DeletedByEach
		}
	} else {
		status = entity.DeletedByTwo
		if friend.Status == entity.DeletedByOne {
			status = entity.DeletedByEach
		}
	}
	return f.friendEntity.UpdateStatus(friend.ID, status)
}

func (f *friendsService) ListFriends(userId int64) ([]FriendInfo, error) {
	friends, err := f.friendEntity.ListFriends(userId)
	if err != nil {
		return nil, err
	}
	friendIds := make([]int64, 0, len(friends))
	for _, item := range friends {
		if item.User1Id != userId {
			friendIds = append(friendIds, item.User1Id)
		} else {
			friendIds = append(friendIds, item.User2Id)
		}
	}
	fs := make([]FriendInfo, 0)
	users, err := f.userEntity.ListUserById(friendIds)
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		fs = append(fs, FriendInfo{
			UserId:   u.ID,
			Username: u.Username,
			Avatar:   u.Avatar,
		})
	}
	return fs, nil
}

func (f *friendsService) IsFriend(userId int64, friendId int64) (bool, error) {
	return f.friendEntity.IsFriend(userId, friendId)
}
