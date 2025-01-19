package wechat

import "github.com/gin-gonic/gin"

type Service interface {
	// Login
	// return unionid openid session_key error
	Login(ctx *gin.Context, code string) (string, string, string, error)
}
