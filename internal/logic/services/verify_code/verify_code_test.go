package verify_code

import (
	"context"
	"github.com/ykds/zura/pkg/cache"
	"testing"
)

var (
	service VerifyCodeService
)

func Init() {
	cfg := cache.DefaultConfig()
	service = NewVerifyCodeService(cache.NewRedis(context.TODO(), &cfg))
}

func Test_verifyCodeService_GenVerifyCode(t *testing.T) {
	Init()
	code, err := service.GenVerifyCode("test")
	if err != nil {
		panic(err)
	}
	t.Log(code)
}

func Test_verifyCodeService_VerifyCode(t *testing.T) {
	Init()
	_, err := service.VerifyCode("test", "123")
	if err != nil {
		panic(err)
	}
}
