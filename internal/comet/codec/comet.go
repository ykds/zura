package codec

import "github.com/ykds/zura/pkg/errors"

const (
	UserIsOffline = 500001 + iota
	MessageIsFull
)

func init() {
	errors.NewCode(UserIsOffline, "用户离线")
	errors.NewCode(MessageIsFull, "消息发送过快")
}
