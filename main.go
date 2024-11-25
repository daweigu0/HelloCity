package HelloCity

import (
	"HelloCity/internal/repository"
	"HelloCity/internal/repository/dao"
	"HelloCity/internal/service"
	"HelloCity/internal/web"
	"HelloCity/ioc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "hellocity.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	userHandler.RegisterRoutes(server)
	server.Run(":8080")
}
