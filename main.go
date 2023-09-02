package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
	"webook/internal/repository"
	entity "webook/internal/repository/entity"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
)

func main() {

	//server := initServer()

	//db := initDB()
	//redisClient := initRedis()
	//
	//userHandler := initUser(db)
	//userHandler.RegisterRoutes(server)
	//// 限流插件：一分钟之内100个请求
	//server.Use(ratelimit.NewBuilder(redisClient, time.Minute, 100).Build())

	server := gin.Default()
	// 先启动服务
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好啊")
		return
	})

	server.Run(":8080")
}

func initServer() *gin.Engine {
	server := gin.Default()

	// 全局中间件
	server.Use(func(ctx *gin.Context) {
		fmt.Println("第一个全局中间件")
	})

	server.Use(func(ctx *gin.Context) {
		fmt.Println("第二个全局中间件")
	})

	// cors中间件
	server.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"https://foo.com"},
		AllowMethods:  []string{"POST"},
		AllowHeaders:  []string{"Content-Type", "authorization"},
		ExposeHeaders: []string{"Content-Length", "X-Jwt-Token"},
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
	}))

	// session中间件
	// 此处表示session存在store中，可以替换成redis
	store := cookie.NewStore([]byte("secret"))
	// 用redis存cookie
	//store := redis.NewStore()
	// 浏览器cookie的key
	server.Use(sessions.Sessions("sessions", store))

	//// 登录鉴权的middleware
	//middleware.IgnorePaths = []string{"/user/login", "/user/signup"}
	//server.Use(middleware.CheckLogin())

	server.Use(middleware.NewLoginJWTMiddlewareBuilder().Ignore("/user/login").Ignore("/user/signup").Build())

	return server
}

// 初始化user模块各个内容
func initUser(db *gorm.DB) *web.UserHandler {
	//	从entity -> repo -> service -> handler
	e := entity.NewUserEntity(db)
	repo := repository.NewUserRepository(e)
	srv := service.NewUserService(repo)
	controller := web.NewUserHandler(srv)
	return controller
}

// 初始化数据库连接
func initDB() *gorm.DB {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:root@tcp(127.0.0.1:13306)/webook"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		//	应该打日志
		// 初始数据库连接失败，server也没必要运行了
		panic(err)
	}

	// 自动迁移表结构
	err = entity.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}

func initRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return redisClient
}
