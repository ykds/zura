package codec

import "zura/pkg/errors"

const (
	StatusErrCode = 300001 + iota
	HadBeFriendCode
	NotFriendCode
	ApplyMySelfErrorCode
	HandleSelfApplyErrCode
	DuplicateHandleApplymentErrCode
)

func init() {
	errors.NewCode(StatusErrCode, "状态错误")
	errors.NewCode(HadBeFriendCode, "已添加该好友")
	errors.NewCode(NotFriendCode, "不是好友关系")
	errors.NewCode(ApplyMySelfErrorCode, "不能给自己提交好友申请")
	errors.NewCode(HandleSelfApplyErrCode, "不能处理自己提交的好友申请")
	errors.NewCode(DuplicateHandleApplymentErrCode, "不能重复处理好友申请")
}
