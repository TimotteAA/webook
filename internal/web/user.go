package web

import (
	"net/http"
	"time"
	"unicode/utf8"
	"webook/internal/domain"
	"webook/internal/service"

	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var RefreshTokenKey = []byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixA")

type UserHandler interface {
	RegisterRoutes(server *gin.Engine)
	Signup(ctx *gin.Context)
	LoginJWT(ctx *gin.Context)
	Signout(ctx *gin.Context)
	Edit(ctx *gin.Context)
	Profile(ctx *gin.Context)
	SignUpCode(ctx *gin.Context)
	LoginByCode(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

// 定义user模块的所有路由
type userHandler struct {
	emailReg    *regexp2.Regexp
	passwordReg *regexp2.Regexp
	srv         service.UserService
	codeService service.CodeService
}

func NewUserHandler(srv service.UserService, codeService service.CodeService) UserHandler {
	// controller入参正则pattern
	const (
		emailRegPattern     = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		passwordRegParttern = `^(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#$%^&*]).{8,60}$`
	)

	emailReg, passwordReg := regexp2.MustCompile(emailRegPattern, regexp2.None), regexp2.MustCompile(passwordRegParttern, regexp2.None)

	u := &userHandler{
		emailReg:    emailReg,
		passwordReg: passwordReg,
		srv:         srv,
		codeService: codeService,
	}
	return u
}

// 统一注册user的路由
func (u *userHandler) RegisterRoutes(server *gin.Engine) {
	//// 统一前缀
	ug := server.Group("/user")
	ug.POST("/signup", u.Signup)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/logout", u.Signout)
	ug.POST("/edit", u.Edit)
	ug.POST("/profile", u.Profile)
	ug.POST("/signup/code/send", u.SignUpCode)
	ug.POST("/login/code", u.LoginByCode)
	ug.POST("/refresh_token", u.RefreshToken)
	//server.POST("/user/login", u.Login)
	//server.POST("/user/logout", u.Signout)
	//server.POST("/user/edit", u.Edit)
}

// 注册路由handler
func (u *userHandler) Signup(ctx *gin.Context) {
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
	// bind出错，gin直接返回
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "Bind data error"})
		return
	}

	isValid, err := u.emailReg.MatchString(req.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	if !isValid {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "邮箱格式不正确"})
		return
	}

	isValid, err = u.passwordReg.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	if !isValid {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "密码必须大于8位，包含数字、特殊字符"})
		return
	}

	if req.Repassword != req.Password {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "密码与确认密码不一致！"})
		return
	}

	// 校验请求入参成功后，调用service方法
	err = u.srv.SignUp(ctx, domain.User{
		Password: req.Password,
		Email:    req.Email,
	})

	if err == service.ErrUserDuplicate {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "邮箱已被注册"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}

	ctx.JSON(http.StatusOK, Result{Code: 0, Msg: "注册成功"})
}

// 登录
func (u *userHandler) LoginJWT(ctx *gin.Context) {
	// 1. 定义请求体
	type ReqUserLogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 2. bind拿结果
	var reqUserLogin *ReqUserLogin

	// 如果在这里打断点能进来，说明中间件没问题
	if err := ctx.Bind(&reqUserLogin); err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}

	// 3. 调用服务方法，注意传值
	user, err := u.srv.Login(ctx, domain.User{Email: reqUserLogin.Email, Password: reqUserLogin.Password})
	if err == service.ErrEmailOrPassWrong {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "邮箱或者密码错误"})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}

	u.setJWTToken(ctx, user)
}

// 退出
func (u *userHandler) Signout(ctx *gin.Context) {
	ctx.String(http.StatusOK, "用户退出")
}

// 编辑
func (u *userHandler) Edit(ctx *gin.Context) {
	type UserEditRequest struct {
		Nickname    string `json:"nickname"`
		Birthday    string `json:"birthday"`
		Description string `json:"description"`
	}

	userEditRequest := &UserEditRequest{}
	err := ctx.Bind(userEditRequest)

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 校验各个字段
	// 可以为空
	nameLen := utf8.RuneCountInString(userEditRequest.Nickname)
	if nameLen > 10 || nameLen < 0 {
		ctx.String(http.StatusBadRequest, "Nickname的长度不能超过10个")
		return
	}

	descLen := utf8.RuneCountInString(userEditRequest.Description)
	if descLen > 300 || descLen < 0 {
		ctx.String(http.StatusBadRequest, "Nickname的长度不能超过300个")
		return
	}

	var birtyTime int64
	if len(userEditRequest.Birthday) > 0 {
		t, err := time.Parse("2006-01-02", userEditRequest.Birthday)
		if err != nil {
			ctx.String(http.StatusBadRequest, "生日日期格式错误")
			return
		}
		birtyTime = t.UnixMilli()
	}

	c, exists := ctx.Get("Claims")
	if !exists {
		ctx.String(http.StatusUnauthorized, "请重新登录")
		return
	}
	claims, ok := c.(*UserJwtClaims)
	if !ok {
		ctx.String(http.StatusUnauthorized, "请重新登录")
		return
	}
	// 调用service
	user, err := u.srv.Edit(ctx, claims.UserId, userEditRequest.Nickname, userEditRequest.Description, birtyTime)

	if err == service.ErrUserNotFound {
		ctx.String(http.StatusBadRequest, "更新的用户不存在")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// 详情
func (u *userHandler) Profile(ctx *gin.Context) {
	var claims *UserJwtClaims

	c, exist := ctx.Get("Claims")
	if !exist {
		ctx.String(http.StatusUnauthorized, "请重新登录")
		return
	}
	claims, exist = c.(*UserJwtClaims)
	if !exist {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	userId := claims.UserId

	user, err := u.srv.FindOne(ctx, userId)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// 发注册短信
func (u *userHandler) SignUpCode(ctx *gin.Context) {
	type signUpBody struct {
		Phone string `json:"phone"`
	}

	var req signUpBody
	// email校验，暂时省略
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	phone := req.Phone

	err := u.codeService.Send(ctx, "login", phone)
	if err == service.ErrCodeSendTooMany {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码发送过于频繁",
		})
		return
	}
	if err == service.ErrUnknownForCode {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "发送成功",
	})
	return
}

func (u *userHandler) LoginByCode(ctx *gin.Context) {
	type userLoginByCode struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}

	var req userLoginByCode
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 可能还需要校验手机号的格式

	// 校验手机验证码
	ok, err := u.codeService.Verify(ctx, "login", req.Phone, req.Code)

	if err == service.ErrUnknownForCode {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码验证次数过多，请重新发送",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码错误，验证失败",
		})
		return
	}
	user, err := u.srv.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	u.setJWTToken(ctx, user)
}

func (u *userHandler) RefreshToken(ctx *gin.Context) {
	// 从header里面拿refresh_token
	refreshTokenStr := ctx.GetHeader("x-jwt-refresh-token")
	if refreshTokenStr == "" {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "刷新token失败"})
		return
	}
	// 保持和jwt中间件中一样的逻辑
	refreshClaims := UserRefreshJwtClaims{}
	refreshToken, err := jwt.ParseWithClaims(refreshTokenStr, &refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return RefreshTokenKey, nil
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 判断payload中的数据，此处认为用火的主键不会为0
	if refreshToken == nil || !refreshToken.Valid || refreshClaims.UserId == int64(0) {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 生成新的token
	u.setJWTToken(ctx, domain.User{Id: refreshClaims.UserId})
}

func (u *userHandler) setJWTToken(ctx *gin.Context, user domain.User) {
	claims := UserJwtClaims{
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
	// 资源请求用这个token
	ctx.Writer.Header().Set("x-jwt-token", tokenStr)

	// 再颁发一个refresh-token
	refreshClaims := UserRefreshJwtClaims{
		UserId: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			// 7天有效期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims)

	// 签发 xxx.xx.xx的字符串
	refreshTokenStr, err := refreshToken.SignedString(RefreshTokenKey)

	// 签发失败
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 把token放到header里去
	ctx.Writer.Header().Set("x-jwt-refresh-token", refreshTokenStr)

	ctx.JSON(http.StatusOK, Result{Code: 0, Msg: "登录成功"})
	return
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

// refresh-token
type UserRefreshJwtClaims struct {
	// 实现Claims接口
	jwt.RegisteredClaims
	// 自己定义的数据
	UserId int64
}
