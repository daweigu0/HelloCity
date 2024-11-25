package main

import (
	"HelloCity/internal/repository"
	"HelloCity/internal/repository/dao"
	"HelloCity/internal/service"
	"HelloCity/internal/web"
	"HelloCity/ioc"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	usDao := dao.NewUserDAO(ioc.InitDB())
	usRepo := repository.NewUserRepositoryHandler(usDao)
	usSvc := service.NewUserServiceHandler(usRepo)
	userHandler := web.NewUserHandler(usSvc)
	userHandler.RegisterRoutes(server)
	server.Run(":8080")
}
