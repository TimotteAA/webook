package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func InitCors() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins:     []string{"https://foo.com"},
		AllowMethods:  []string{"POST"},
		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length", "X-Jwt-Token", "expire-time"},
		// 允许携带cookie
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			// 开发环境
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			// 公司域名
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	})
}
