package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
	"webook/internal/web"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

// Builder模式
func (m *LoginJWTMiddlewareBuilder) Ignore(path string) *LoginJWTMiddlewareBuilder {
	m.paths = append(m.paths, path)
	return m
}

// 最后真正的中间件
// 校验jwt token是否存在
func (m *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 对指定路径无视中间件
		for _, path := range m.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 判断header
		tokenHeader := ctx.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			//	没带token
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 通常是 Bearer xxx
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 真正的tokenString: xxxx.xxxx.xxxx
		// 校验jwt token
		tokenStr := segs[1]
		claims := &web.UserJwtClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 判断jwt payload，此处等于0也许不行
		if token == nil || !token.Valid || claims.UserId == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//	判断UserAgent
		if ctx.Request.UserAgent() != claims.UserAgent {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// token续期
		now := time.Now()
		// 距离签发时间过了10s，就刷新
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			// 重新设定过期时间
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			// 办法新token
			tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
			if err != nil {
				// 刷新失败
				log.Println("jwt刷新失败 ", err)
			}
			fmt.Println("签发失败吗 ", err)
			ctx.Header("x-jwt-token", tokenStr)
		}
		fmt.Println("jwt中间件 ", claims)
		// 拿到user信息，供后面的路由使用，注意这里放的是指针，断言应该也是指针
		ctx.Set("Claims", claims)
	}
}
