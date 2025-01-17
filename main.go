package main

import (
	_ "HelloCity/docs"
	"HelloCity/internal/repository"
	"HelloCity/internal/repository/cache"
	"HelloCity/internal/repository/dao"
	"HelloCity/internal/service"
	"HelloCity/internal/service/oss/qiniu"
	"HelloCity/internal/utils"
	"HelloCity/internal/web"
	"HelloCity/ioc"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"strings"
	"time"
)

func InitMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
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
		}),
		//(&middleware.LoginJWTMiddlewareBuilder{}).CheckLogin(),
	}
}

func InitWebServer() *gin.Engine {
	server := gin.Default()
	server.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	usDao := dao.NewUserDAO(ioc.InitDB())
	usRepo := repository.NewUserRepositoryHandler(usDao)
	usSvc := service.NewUserService(usRepo)
	tkCache := cache.NewTokenCache(ioc.InitRedis(), time.Minute*5)
	tkRepo := repository.NewTokenCachedRepository(tkCache)
	tkSvc := service.NewTokenService(tkRepo)
	userHandler := web.NewUserHandler(usSvc, tkSvc)
	server.Use(InitMiddlewares()...)
	userHandler.RegisterRoutes(server)
	accessKey := utils.Config.GetString("oss.qiniu.accessKey")
	secretKey := utils.Config.GetString("oss.qiniu.secretKey")
	fileSvc := qiniu.NewService(accessKey, secretKey)
	fileHandler := web.NewFileHandler(fileSvc)
	fileHandler.RegisterRoutes(server)
	return server
}

// @title 你好同城后端API
// @accept json
// @produce	json
// @schemes	http https
func main() {
	server := InitWebServer()
	domain := utils.Config.GetString("nihaotongcheng.domain")
	port := utils.Config.GetString("nihaotongcheng.port")
	err := server.Run(fmt.Sprintf("%s:%s", domain, port))
	if err != nil {
		log.Panicf("服务器启动错误 %v\n", err)
	}
}
