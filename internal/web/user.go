package web

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"webook/internal/domain"
	"webook/internal/service"
)

// 定义user模块的所有路由
type UserHandler struct {
	emailReg    *regexp2.Regexp
	passwordReg *regexp2.Regexp
	srv         *service.UserService
}

func NewUserHandler(srv *service.UserService) *UserHandler {
	// controller入参正则pattern
	const (
		emailRegPattern     = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		passwordRegParttern = `^(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#$%^&*]).{8,60}$`
	)

	emailReg, passwordReg := regexp2.MustCompile(emailRegPattern, regexp2.None), regexp2.MustCompile(passwordRegParttern, regexp2.None)

	u := &UserHandler{
		emailReg:    emailReg,
		passwordReg: passwordReg,
		srv:         srv,
	}
	return u
}

// 统一注册user的路由
func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	//// 统一前缀
	ug := server.Group("/user")
	ug.POST("/signup", u.Signup)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/logout", u.Signout)
	ug.POST("/edit", u.Edit)

	//server.POST("/user/login", u.Login)
	//server.POST("/user/logout", u.Signout)
	//server.POST("/user/edit", u.Edit)
}

// 注册路由handler
func (u *UserHandler) Signup(ctx *gin.Context) {
	// 注册请求结构体
	type SignUpReq struct {
		// 此处后面的 json:"xxx"，表示从body里的某个字段取数据
		Email      string `json:"email"`
		Password   string `json:"password"`
		Repassword string `json:"repassword"`
	}

	// 根据content-type序列化body数据
	var req *SignUpReq
	// 注意穿的是值，而不是指针
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusBadRequest, "请求参数错误")
		return
	}
	fmt.Println("请求入参 ", req)

	isValid, err := u.emailReg.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isValid {
		ctx.String(http.StatusOK, "邮箱格式不正确")
		return
	}

	isValid, err = u.passwordReg.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isValid {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字、特殊字符")
		return
	}

	if req.Repassword != req.Password {
		ctx.String(http.StatusBadRequest, "密码与确认密码不一致！")
		return
	}

	// 校验请求入参成功后，调用service方法
	err = u.srv.SignUp(ctx, domain.User{
		Password: req.Password,
		Email:    req.Email,
	})

	if err == service.ErrUserDuplciateEmail {
		ctx.String(http.StatusOK, "邮箱重复")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "用户注册")
}

// 登录
func (u *UserHandler) Login(ctx *gin.Context) {
	// 1. 定义请求体
	type ReqUserLogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 2. bind拿结果
	var reqUserLogin *ReqUserLogin

	// 如果在这里打断点能进来，说明中间件没问题
	if err := ctx.Bind(&reqUserLogin); err != nil {

		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 3. 调用服务方法，注意传值
	user, err := u.srv.Login(ctx, domain.User{Email: reqUserLogin.Email, Password: reqUserLogin.Password})
	if err == service.ErrEmailOrPassWrong {
		ctx.String(http.StatusOK, "邮箱或者密码错误")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 设置session
	session := sessions.Default(ctx)
	session.Set("userId", user.Id)
	session.Options(sessions.Options{
		Secure:   true,
		HttpOnly: true,
		// 一分钟过期
		MaxAge: 600,
	})
	session.Save()
	ctx.String(http.StatusOK, "登录成功")
	return
}

// 登录
func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	// 1. 定义请求体
	type ReqUserLogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 2. bind拿结果
	var reqUserLogin *ReqUserLogin

	// 如果在这里打断点能进来，说明中间件没问题
	if err := ctx.Bind(&reqUserLogin); err != nil {

		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 3. 调用服务方法，注意传值
	user, err := u.srv.Login(ctx, domain.User{Email: reqUserLogin.Email, Password: reqUserLogin.Password})
	if err == service.ErrEmailOrPassWrong {
		ctx.String(http.StatusOK, "邮箱或者密码错误")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserJwtClaims{
		UserId:    user.Id,
		UserAgent: ctx.Request.UserAgent(),
	})

	// 签发 xxx.xx.xx的字符串
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))

	// 签发失败
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 把token放到header里去
	ctx.Writer.Header().Set("x-jwt-token", tokenStr)

	ctx.String(http.StatusOK, "登录成功")
	return
}

// 退出
func (u *UserHandler) Signout(ctx *gin.Context) {
	ctx.String(http.StatusOK, "用户退出")
}

// 编辑
func (u *UserHandler) Edit(ctx *gin.Context) {
	ctx.String(http.StatusOK, "编辑用户")
}

// 自定义jwt-claims
type UserJwtClaims struct {
	// 实现Claims接口
	jwt.RegisteredClaims
	// 自己定义的数据
	UserId int64
	// 浏览器信息
	UserAgent string
}
