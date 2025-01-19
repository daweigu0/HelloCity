//go:build wireinject

package main

import (
	"HelloCity/internal/repository"
	"HelloCity/internal/repository/cache"
	"HelloCity/internal/repository/dao"
	"HelloCity/internal/service"
	"HelloCity/internal/service/wechat/power_wechat"
	"HelloCity/internal/web"
	"HelloCity/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	cache.NewTokenCache,
	cache.NewUserCache,
	dao.NewUserDAO,
	repository.NewTokenCachedRepository,
	repository.NewUserRepositoryHandler,
	ioc.NewOssService,
	service.NewTokenService,
	service.NewUserService,
	web.NewUserHandler,
	web.NewFileHandler,
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitWebServer,
	ioc.NewTimeDuration,
	ioc.NewWechatService,
	power_wechat.NewService,
)

func InitWebServer() *gin.Engine {
	wire.Build(ProviderSet)
	return nil
}
