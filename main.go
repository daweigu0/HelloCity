package main

import (
	_ "HelloCity/docs"
	"HelloCity/internal/repository"
	"HelloCity/internal/repository/dao"
	"HelloCity/internal/service"
	"HelloCity/internal/web"
	"HelloCity/internal/web/middleware"
	"HelloCity/ioc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
		(&middleware.LoginJWTMiddlewareBuilder{}).CheckLogin(),
	}
}

func InitWebServer() *gin.Engine {
	server := gin.Default()
	server.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	usDao := dao.NewUserDAO(ioc.InitDB())
	usRepo := repository.NewUserRepositoryHandler(usDao)
	usSvc := service.NewUserServiceHandler(usRepo)
	userHandler := web.NewUserHandler(usSvc)
	server.Use(InitMiddlewares()...)
	userHandler.RegisterRoutes(server)

	return server
}

// @title 你好同城后端API
// @accept json
// @produce	json
// @schemes	http https
func main() {
	server := InitWebServer()
	server.Run(":8080")
}
