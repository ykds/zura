package errors

import (
	"log"
)

var codeMap = map[int]Code{}

var (
	OKStatus             = 200
	ParameterErrorStatus = 400
	UnAuthorizedStatus   = 401
	InternalErrorStatus  = 500
)

func init() {
	NewCode(OKStatus, "成功")
	NewCode(ParameterErrorStatus, "参数错误")
	NewCode(UnAuthorizedStatus, "未认证")
	NewCode(InternalErrorStatus, "服务器错误")
}

type Code struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewCode(status int, message string) Code {
	if _, ok := codeMap[status]; ok {
		log.Panicf("code %v already exists", status)
	}
	c := Code{
		Status:  status,
		Message: message,
	}
	codeMap[status] = c
	return c
}

func GetCode(status int) Code {
	if c, ok := codeMap[status]; ok {
		return c
	}
	log.Printf("Code(%d)未定义", status)
	return codeMap[InternalErrorStatus]
}
