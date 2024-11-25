package HelloCity

import (
	"HelloCity/internal/web"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	userHandler := web.NewUserHandler()
	userHandler.RegisterRoutes(server)
	server.Run(":6666")
}
