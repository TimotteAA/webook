package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
	"webook/internal/service"
	"webook/internal/web/ijwt"
)

var stateJWTKey = []byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixB")
var errStateWrong = errors.New("state被篡改了")

type OAuth2WeChatHandler interface {
	RegisterRoutes(server *gin.Engine)
	AuthUrl(ctx *gin.Context)
	Callback(ctx *gin.Context)
}

type oAuth2WeChatHandler struct {
	svc         service.WeChatService
	userService service.UserService
}

func NewOAuth2WeChatHandler(svc service.WeChatService) OAuth2WeChatHandler {
	return &oAuth2WeChatHandler{svc: svc}
}

func (handler *oAuth2WeChatHandler) RegisterRoutes(server *gin.Engine) {
	group := server.Group("/oauth2/wechat")
	// 申请登录的第三方地址
	group.GET("/authurl", handler.AuthUrl)
	// 扫码后的callback url，用Any保险一点
	group.Any("/callback", handler.Callback)
}

func (handler *oAuth2WeChatHandler) AuthUrl(ctx *gin.Context) {
	// 标识此次扫码的code，类似于单次会话id
	state := uuid.New()
	url, err := handler.svc.AuthURL(state.String())
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	// 设置cookie
	err = handler.setStateCookie(ctx, state.String())
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}

	ctx.JSON(http.StatusOK, Result{Msg: "success", Code: 0, Data: url})
	return
}

func (handler *oAuth2WeChatHandler) Callback(ctx *gin.Context) {
	err := handler.verifyStateCookie(ctx)
	if err == errStateWrong {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "登录失效，请重新登录"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	code := ctx.Query("code")
	result, err := handler.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	user, err := handler.userService.FindOrCreateByWeChat(ctx, result)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 颁发token: todo封装方法
	claims := ijwt.UserJwtClaims{
		UserId:    user.Id,
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 600)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// 签发 xxx.xx.xx的字符串
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))

	// 签发失败
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 把token放到header里去
	ctx.Writer.Header().Set("x-jwt-token", tokenStr)
	ctx.JSON(http.StatusOK, Result{Code: 0, Msg: "登陆成功"})
	return
}

func (handler *oAuth2WeChatHandler) setStateCookie(ctx *gin.Context, state string) error {
	claims := WeChatStateClaims{
		State: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(stateJWTKey)
	if err != nil {
		return err
	}
	// 放在cookie里，再callback时进行校验
	ctx.SetCookie("ijwt-oauth2-state", tokenStr, 600, "/oauth2/wechat/callback", "", false, true)
	return nil
}

func (handler *oAuth2WeChatHandler) verifyStateCookie(ctx *gin.Context) error {
	// url上的state
	state := ctx.Query("state")
	tokenStr, err := ctx.Cookie("ijwt-oauth2-state")
	// 拿不到cookie
	if err != nil || tokenStr == "" {
		return fmt.Errorf("%w, 无法拿到cookie", err)
	}
	// 校验tokenStr
	claims := &WeChatStateClaims{}
	_, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return stateJWTKey, nil
	})
	if err != nil {
		return fmt.Errorf("%w, token解析失败", err)
	}
	if state != claims.State {
		return errStateWrong
	}
	return nil
}

type WeChatStateClaims struct {
	jwt.Claims
	State string
}
