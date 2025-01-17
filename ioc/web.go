package ioc

import (
	"HelloCity/internal/web"
	"HelloCity/internal/web/middleware"
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

func InitWebServer(userHdl *web.UserHandler, fileHdl *web.FileHandler) *gin.Engine {
	server := gin.Default()
	server.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.Use(InitMiddlewares()...)
	userHdl.RegisterRoutes(server)
	fileHdl.RegisterRoutes(server)
	return server
}
