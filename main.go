package main

import (
	"HelloCity/internal/repository"
	"HelloCity/internal/repository/dao"
	"HelloCity/internal/service"
	"HelloCity/internal/web"
	"HelloCity/ioc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	server := gin.Default()
	usDao := dao.NewUserDAO(ioc.InitDB())
	usRepo := repository.NewUserRepositoryHandler(usDao)
	usSvc := service.NewUserServiceHandler(usRepo)
	userHandler := web.NewUserHandler(usSvc)
	server.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "nihaotongcheng.com")
		},
		MaxAge: 12 * time.Hour,
	}), func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/login" {
			// 不需要登录校验
			return
		}
		// 根据约定，token 在 Authorization 头部
		// Bearer XXXX
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(authCode, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		var uc web.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if uc.UserAgent != ctx.GetHeader("User-Agent") {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		expireTime := uc.ExpiresAt
		if expireTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 5))
			tokenStr, err = token.SignedString(web.JWTKey)
			ctx.Header("x-jwt-token", tokenStr)
			if err != nil {
				// 这边不要中断，因为仅仅是过期时间没有刷新，但是用户是登录了的
				log.Println(err)
			}
		}
		ctx.Set("user", uc)
	})

	userHandler.RegisterRoutes(server)
	server.Run(":8080")
}
