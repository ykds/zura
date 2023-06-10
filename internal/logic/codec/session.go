package codec

import "zura/pkg/errors"

const (
	AddMemberToPointSessionErrCode = iota + 400001
	NotPermitChangeRole
)

func init() {
	errors.NewCode(AddMemberToPointSessionErrCode, "私聊会话不能添加成员")
	errors.NewCode(NotPermitChangeRole, "无权限分配角色")
}
