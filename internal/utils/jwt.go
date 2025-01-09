package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgA")

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       uint64
	UserAgent string
}

func GenerateToken(uc *UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		return "", fmt.Errorf("jwt生成token错误:%v", err)
	}
	return tokenString, nil
}
