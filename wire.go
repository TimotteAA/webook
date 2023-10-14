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
	"webook/internal/web/ijwt"
	"webook/ioc"
	"webook/ioc/oauth2"
	"webook/ioc/sms"
)

// 根目录下运行wire
func InitWebServer() *gin.Engine {
	wire.Build(

		// ijwt handler
		ijwt.NewRedisHandler,

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
		oauth2.InitWeChatService,

		// service
		//local.NewMemoryService,
		sms.InitSmsService,
		service.NewCodeService,
		service.NewUserService,

		// controller
		web.NewUserHandler,
		web.NewOAuth2WeChatHandler,

		ioc.InitMiddlewares,
		ioc.InitWebServer,
	)
	return new(gin.Engine)
}
