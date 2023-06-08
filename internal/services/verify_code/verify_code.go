package verify_code

import (
	"context"
	"time"
	"zura/pkg/cache"
)

func NewVerifyCodeService(cache *cache.Redis) VerifyCodeService {
	return &verifyCodeService{
		cache: cache,
	}
}

type VerifyCodeService interface {
	GenVerifyCode(key string) error
	VerifyCode(key string, code string) (bool, error)
}

type verifyCodeService struct {
	cache *cache.Redis
}

func (v *verifyCodeService) GenVerifyCode(key string) error {
	panic("not implemented") // TODO: Implement
}

func (v *verifyCodeService) VerifyCode(key string, code string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	verifyCode, err := v.cache.Get(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return code == verifyCode, nil
}
