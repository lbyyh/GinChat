package tools

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	jwt.StandardClaims
}

var jwtSecret = []byte("your_secret_key")

// GetJwt 生成JWT Token
func GetJwt(userID int64, name string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		UserID: userID,
		Name:   name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "ginchat",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}
