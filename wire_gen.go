// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"HelloCity/internal/repository"
	"HelloCity/internal/repository/cache"
	"HelloCity/internal/repository/dao"
	"HelloCity/internal/service"
	"HelloCity/internal/web"
	"HelloCity/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

import (
	_ "HelloCity/docs"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	db := ioc.InitDB()
	userDao := dao.NewUserDAO(db)
	userRepository := repository.NewUserRepositoryHandler(userDao)
	userService := service.NewUserService(userRepository)
	cmdable := ioc.InitRedis()
	duration := ioc.NewTimeDuration()
	tokenCache := cache.NewTokenCache(cmdable, duration)
	tokenRepository := repository.NewTokenCachedRepository(tokenCache)
	tokenService := service.NewTokenService(tokenRepository)
	ossService := ioc.NewOssService()
	userHandler := web.NewUserHandler(userService, tokenService, ossService)
	fileHandler := web.NewFileHandler(ossService)
	engine := ioc.InitWebServer(userHandler, fileHandler)
	return engine
}

// wire.go:

var ProviderSet = wire.NewSet(cache.NewTokenCache, cache.NewUserCache, dao.NewUserDAO, repository.NewTokenCachedRepository, repository.NewUserRepositoryHandler, ioc.NewOssService, service.NewTokenService, service.NewUserService, web.NewUserHandler, web.NewFileHandler, ioc.InitDB, ioc.InitRedis, ioc.InitWebServer, ioc.NewTimeDuration)
