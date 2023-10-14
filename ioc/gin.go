package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/internal/web"
	"webook/internal/web/ijwt"
	"webook/internal/web/middleware"
	"webook/pkg/ginx/middlewares/ratelimit"
)

// 初始gin的服务器
func InitWebServer(ug web.UserHandler, wechatHandler web.OAuth2WeChatHandler, fn []gin.HandlerFunc) *gin.Engine {
	server := gin.Default()
	server.Use(fn...)
	ug.RegisterRoutes(server)
	wechatHandler.RegisterRoutes(server)
	return server
}

func InitMiddlewares(cmd redis.Cmdable, handler ijwt.Handler) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.InitCors(),
		ratelimit.NewBuilder(cmd, time.Minute, 100).Build(),
		middleware.NewLoginJWTMiddlewareBuilder(cmd, handler).
			Ignore("/user/login").
			Ignore("/user/signup").
			Ignore("/user/signup/code/send").
			Ignore("/user/login/code").
			Ignore("/user/refresh_token").
			Build(),
	}
}
