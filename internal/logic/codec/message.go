package codec

import "github.com/ykds/zura/pkg/errors"

const (
	IllegalMsgTsCode = iota + 600001
)

func init() {
	errors.NewCode(IllegalMsgTsCode, "非法消息时间戳")
}
