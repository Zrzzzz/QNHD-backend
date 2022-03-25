package util

import (
	"errors"
	"fmt"
	"qnhd/pkg/setting"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtSecret       = []byte(setting.AppSetting.JwtSecret)
	ErrInvalidToken = fmt.Errorf("invalid token")
)

type Claims struct {
	Uid string `json:"uid"`
	jwt.StandardClaims
}

func GenerateToken(uid string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(setting.AppSetting.TokenExpireTime) * time.Hour)
	claims := Claims{
		uid,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "qnhd",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return jwtSecret, nil
	})

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrInvalidToken) {
			return nil, ErrInvalidToken
		}
		return nil, ErrInvalidToken
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
