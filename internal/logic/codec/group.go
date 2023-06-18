package codec

import "github.com/ykds/zura/pkg/errors"

const (
	NotPermitCode = iota + 500001
	NotGroupMember
	HadAddGroupCode
	UnSupportRoleCode
)

func init() {
	errors.NewCode(NotPermitCode, "无权限")
	errors.NewCode(NotGroupMember, "非该群成员")
	errors.NewCode(HadAddGroupCode, "已是该群成员")
	errors.NewCode(UnSupportRoleCode, "不支持该角色")
}
