package verify_code

import (
	"context"
	"github.com/ykds/zura/pkg/cache"
	"github.com/ykds/zura/pkg/errors"
	"github.com/ykds/zura/pkg/random"
	"time"
)

func NewVerifyCodeService(cache cache.Cache) VerifyCodeService {
	return &verifyCodeService{
		cache: cache,
	}
}

type VerifyCodeService interface {
	GenVerifyCode(key string) (string, error)
	VerifyCode(key string, code string) (bool, error)
}

type verifyCodeService struct {
	cache cache.Cache
}

func (v *verifyCodeService) GenVerifyCode(key string) (string, error) {
	code := random.RandNum(8)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := v.cache.Set(ctx, key, code, time.Minute*30)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (v *verifyCodeService) VerifyCode(key string, code string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	verifyCode, err := v.cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, cache.NotFoundErr) {
			return false, nil
		}
		return false, err
	}
	return code == verifyCode.(string), nil
}
