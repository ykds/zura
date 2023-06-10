package codec

import "zura/pkg/errors"

const (
	StatusErrCode = 300001 + iota
	HadBeFriendCode
	NotFriendCode
)

func init() {
	errors.NewCode(StatusErrCode, "状态错误")
	errors.NewCode(HadBeFriendCode, "已添加该好友")
	errors.NewCode(NotFriendCode, "不是好友关系")
}
