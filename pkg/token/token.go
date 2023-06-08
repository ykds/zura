package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UserId int64 `json:"user_id"`
}

func NewToken(userId int64) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * 7)),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserId: userId,
	})
	return token.SignedString([]byte("zura"))
}

func VerifyToken(token string) (int64, error) {
	t, err := jwt.ParseWithClaims(token, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected method")
		}
		return []byte("zura"), nil
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}
	if !t.Valid {
		return 0, errors.New("token已失效")
	}
	if c, ok := t.Claims.(*UserClaims); ok {
		return c.UserId, nil
	}
	return 0, errors.New("token已失效")
}
