//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/entity"
	"webook/internal/service"
	"webook/internal/web"
	"webook/ioc"
	"webook/ioc/sms"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// db和redis
		ioc.InitDB,
		ioc.InitRedis,

		// cache和entity
		entity.NewUserEntity,
		cache.NewCodeCache,
		cache.NewUserCache,

		// repo
		repository.NewUserRepository,
		repository.NewCodeRepository,

		// wechat service

		// service
		//local.NewMemoryService,
		sms.InitSmsService,
		service.NewCodeService,
		service.NewUserService,

		// controller
		web.NewUserHandler,

		ioc.InitMiddlewares,
		ioc.InitWebServer,
	)
	return new(gin.Engine)
}
