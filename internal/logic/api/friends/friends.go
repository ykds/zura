package friends

import (
	"zura/internal/logic/common"
	"zura/internal/logic/services"
	"zura/internal/logic/services/friends"
	"zura/pkg/errors"
	"zura/pkg/response"

	"github.com/gin-gonic/gin"
)

func ListFriends(c *gin.Context) {
	var (
		err  error
		resp struct {
			Data []friends.FriendInfo `json:"data"`
		}
	)
	defer func() {
		if len(resp.Data) == 0 {
			resp.Data = []friends.FriendInfo{}
		}
		response.HttpResponse(c, err, resp)
	}()
	userId := c.GetInt64(common.UserIdKey)
	resp.Data, err = services.GetServices().FriendsService.ListFriends(userId)
}

func DeleteFriends(c *gin.Context) {
	var (
		err error
		req struct {
			FriendId int64 `json:"friend_id"`
		}
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	if err = c.Bind(&req); err != nil {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	if req.FriendId == 0 {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	userId := c.GetInt64(common.UserIdKey)
	err = services.GetServices().FriendsService.DeleteFriend(userId, req.FriendId)
}
