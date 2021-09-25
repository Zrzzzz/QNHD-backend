package util

import (
	"qnhd/pkg/setting"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(setting.AppSetting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Tag      int    `json:"tag"`
	jwt.StandardClaims
}

const (
	ADMIN = 0x53
	USER  = 0x16
)

func GenerateToken(username string, role int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)
	tag := USER
	if role != 1 {
		tag = ADMIN
	}
	claims := Claims{
		username,
		tag,
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
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
