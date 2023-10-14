package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"webook/internal/web/ijwt"
)

/*
*
jwt相关的点
1. claims中存储哪些用户信息？暂时只考虑userId
2. token是无状态的，根本不知道自己啥情况，如何提高安全性？此处在claims中考虑了UserAgent
3. 长短token，如何续期？前端做？后端自己续
4. 如何让token过期？是存储token？此处在claims中加入一个ssid，在redis中进行记录（其实这么做跟session没区别了）
5. redis挂了如何降级？
*/
type LoginJWTMiddlewareBuilder struct {
	paths   []string
	client  redis.Cmdable
	handler ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(client redis.Cmdable, handler ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{client: client, handler: handler}
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

		//// 判断header
		//tokenHeader := ctx.Request.Header.Get("Authorization")
		//if tokenHeader == "" {
		//	//	没带token
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		//// 通常是 Bearer xxx
		//segs := strings.Split(tokenHeader, " ")
		//if len(segs) != 2 {
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		//
		//// 真正的tokenString: xxxx.xxxx.xxxx
		//// 校验jwt token
		//tokenStr := segs[1]

		tokenStr := m.handler.ExtractTokenString(ctx)
		if tokenStr == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := &ijwt.UserJwtClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		})
		if err != nil {
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
		// 判断是否已经退出登录
		logout, err := m.client.Exists(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid)).Result()
		// redis崩溃、或者用户退出
		// 也许可以扩展的点，如果redis崩溃了，就不校验ssid了，作为一个降级策略
		// 如果redis没有崩溃，则校验ssid
		if err != nil {
			// log一下
			log.Println("记录用户ssid的中间件挂了")
		} else if logout > 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//// token续期
		//now := time.Now()
		//// 距离签发时间过了10s，就刷新
		//if claims.ExpiresAt.Sub(now) < time.Second*50 {
		//	// 重新设定过期时间
		//	claims.ExpiresAt = ijwt.NewNumericDate(time.Now().Add(time.Minute))
		//	// 办法新token
		//	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
		//	if err != nil {
		//		// 刷新失败
		//		log.Println("jwt刷新失败 ", err)
		//	}
		//	fmt.Println("签发失败吗 ", err)
		//	ctx.Header("x-ijwt-token", tokenStr)
		//}

		// 拿到user信息，供后面的路由使用，注意这里放的是指针，断言应该也是指针
		ctx.Set("Claims", claims)
	}
}
