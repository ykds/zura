package friends

import (
	"github.com/gin-gonic/gin"
	"github.com/ykds/zura/internal/common"
	"github.com/ykds/zura/internal/logic/services"
	"github.com/ykds/zura/internal/logic/services/friends"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/response"
	"strconv"
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
	if err == nil {
		for i, u := range resp.Data {
			resp.Data[i].Avatar = common.ParseAvatarUrl(c, u.Avatar)
		}
	}
}

func DeleteFriends(c *gin.Context) {
	var (
		err error
	)
	defer func() {
		response.HttpResponse(c, err, nil)
	}()
	s := c.Param("id")
	if s == "" {
		err = errors.WithMessage(errors.New(errors.ParameterErrorStatus), err.Error())
		return
	}
	var id int64
	id, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		return
	}
	userId := c.GetInt64(common.UserIdKey)
	err = services.GetServices().FriendsService.DeleteFriend(userId, id)
}
