package codec

import "github.com/ykds/zura/pkg/errors"

const (
	HeartBeatFailedCode = iota + 700001
	SyncNewMessageFailedCode
)

func init() {
	errors.NewCode(HeartBeatFailedCode, "心跳失败，连接断开")
	errors.NewCode(SyncNewMessageFailedCode, "拉取新消息失败")
}
