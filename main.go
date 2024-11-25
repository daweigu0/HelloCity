package main

import (
	"HelloCity/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func main() {
	server := gin.Default()
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
	userHandler := web.NewUserHandler()
	userHandler.RegisterRoutes(server)
	server.Run(":8888")
}
