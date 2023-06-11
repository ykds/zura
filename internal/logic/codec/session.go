package codec

import "zura/pkg/errors"

const (
	AddMemberToPointSessionErrCode = iota + 400001
	NotPermitChangeRole
	NotPermitDismissGroupCode
	OpenWithSelfErrCode
)

func init() {
	errors.NewCode(AddMemberToPointSessionErrCode, "私聊会话不能添加成员")
	errors.NewCode(NotPermitChangeRole, "无权限分配角色")
	errors.NewCode(OpenWithSelfErrCode, "不能与自己创建会话")
	errors.NewCode(NotPermitDismissGroupCode, "只有群主才能解散群")
}
