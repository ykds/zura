package codec

import "github.com/ykds/zura/pkg/errors"

const (
	OpenWithSelfErrCode = iota + 400001
	UnSupportSessionType
)

func init() {
	errors.NewCode(OpenWithSelfErrCode, "不能与自己创建会话")
	errors.NewCode(UnSupportSessionType, "不支持该会话类型")
}
