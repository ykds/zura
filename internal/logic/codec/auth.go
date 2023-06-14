package codec

import "github.com/ykds/zura/pkg/errors"

const (
	NeedAuthStatus = iota + 100001
	ParseTokenFailedStatus
)

func init() {
	errors.NewCode(NeedAuthStatus, "token为空")
	errors.NewCode(ParseTokenFailedStatus, "解析token失败")
}
