package ijwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

var AccessTokenKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
var RefreshTokenKey = []byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixA")

var ErrUserLogout = errors.New("用户已经退出登录")

// 内聚所有Jwt token相关的方法
type Handler interface {
	ClearToken(ctx *gin.Context) error
	SetLoginToken(ctx *gin.Context, userId int64) error
	SetAccessToken(ctx *gin.Context, userId int64, ssid string) error
	SetRefreshToken(ctx *gin.Context, userId int64, ssid string) error
	CheckSession(ctx *gin.Context, ssid string) error
	ExtractTokenString(ctx *gin.Context) string
}

type RedisHandler struct {
	client redis.Cmdable
	// 长token过期时间
	rtExpiration time.Duration
}

func (r *RedisHandler) ClearToken(ctx *gin.Context) error {
	//	所有的请求都必须有这讲个token
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-jwt-refresh-token", "")
	// 这里不可能拿不到
	claims := ctx.MustGet("Claims").(*UserJwtClaims)
	return r.client.Set(ctx, r.key(claims.Ssid), claims.Ssid, r.rtExpiration).Err()
}

// 设置长短token
func (r *RedisHandler) SetLoginToken(ctx *gin.Context, userId int64) error {
	ssid := uuid.New().String()
	// 先设置短token
	if err := r.SetAccessToken(ctx, userId, ssid); err != nil {
		return err
	}
	// 然后设置长token
	return r.SetRefreshToken(ctx, userId, ssid)
}

func (r *RedisHandler) SetAccessToken(ctx *gin.Context, userId int64, ssid string) error {
	claims := UserJwtClaims{
		UserId:    userId,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 7)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// 签发 xxx.xx.xx的字符串
	tokenStr, err := token.SignedString(AccessTokenKey)
	if err != nil {
		return err
	}
	ctx.Header("x-ijwt-token", tokenStr)
	return nil
}

func (r *RedisHandler) SetRefreshToken(ctx *gin.Context, userId int64, ssid string) error {
	// 再颁发一个refresh-token
	refreshClaims := UserRefreshJwtClaims{
		UserId: userId,
		Ssid:   ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			// 7天有效期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims)

	// 签发 xxx.xx.xx的字符串
	refreshTokenStr, err := refreshToken.SignedString(RefreshTokenKey)
	if err != nil {
		return err
	}
	ctx.Header("x-ijwt-refresh-token", refreshTokenStr)
	return nil
}

// redis挂了校验？
func (r *RedisHandler) CheckSession(ctx *gin.Context, ssid string) error {
	logout, err := r.client.Exists(ctx, r.key(ssid)).Result()
	if err != nil {
		return err
	}
	if logout > 0 {
		return ErrUserLogout
	}
	return nil
}

func (r *RedisHandler) ExtractTokenString(ctx *gin.Context) string {
	// 判断header
	tokenHeader := ctx.Request.Header.Get("Authorization")
	if tokenHeader == "" {
		//	没带token
		return ""
	}
	// 通常是 Bearer xxx
	segs := strings.Split(tokenHeader, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

func (r *RedisHandler) key(ssid string) string {
	return fmt.Sprintf("users:ssid:%s", ssid)
}

func NewRedisHandler(client redis.Cmdable) Handler {
	return &RedisHandler{client: client, rtExpiration: time.Hour * 24 * 7}
}

// 自定义jwt-claims
type UserJwtClaims struct {
	// 实现Claims接口
	jwt.RegisteredClaims
	// 用于标识是否过期
	Ssid string
	// 自己定义的数据
	UserId int64
	// 浏览器信息
	UserAgent string
}

// refresh-token
type UserRefreshJwtClaims struct {
	// 实现Claims接口
	jwt.RegisteredClaims
	// 用于标识是否过期
	Ssid string
	// 自己定义的数据
	UserId int64
}
