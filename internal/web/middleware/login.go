package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

var IgnorePaths []string

// 对指定路由校验session
func CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 注册、登录接口不需要校验session
		for _, path := range IgnorePaths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		session := sessions.Default(ctx)
		id := session.Get("userId")
		if id == nil {
			// 没有登陆
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
