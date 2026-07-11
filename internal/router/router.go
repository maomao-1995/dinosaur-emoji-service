package router

import (
	"dinosaur-emoji-service/internal/handler"
	"dinosaur-emoji-service/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 全局跨域配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有域名，生产可指定前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "refresh_token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 配置静态文件服务
	r.Static("/uploads", "./uploads")

	// Logger中间件
	r.Use(middleware.Logger())

	// 注册Swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务
	// 注意：默认是阻塞式的，会一直运行直到被中断

	//global路由
	r.POST("/sendCode", handler.SendCode)
	r.GET("/refresh", handler.Refresh)
	r.POST("/upload", handler.Upload)

	// user路由
	userGroup := r.Group("/user")
	userPrivate := userGroup.Group("")
	userPrivate.Use(middleware.ParseToken())
	{
		userPrivate.GET("/info", handler.UserInfo)
	}
	userPublic := userGroup.Group("")
	{
		userPublic.POST("/register", handler.UserRegister)
		userPublic.POST("/login", handler.Login)
		userPublic.POST("/loginAndRegister", handler.LoginAndRegister)

	}

	//emoji路由
	emojiGroup := r.Group("/emoji")
	emojiGroup.Use(middleware.ParseToken())
	{
		emojiGroup.GET("/detail", handler.EmojiDetail)
		emojiGroup.POST("/add", handler.EmojiAdd)
		emojiGroup.POST("/delete", handler.EmojiDelete)
		emojiGroup.POST("/edit", handler.EmojiEdit)
	}

	//emojiPack路由
	emojiPackGroup := r.Group("/emojiPack")
	emojiPackPrivata := emojiPackGroup.Group("")
	emojiPackPrivata.Use(middleware.ParseToken())
	{
		emojiPackPrivata.POST("/add", handler.EmojiPackAdd)
		emojiPackPrivata.POST("/edit", handler.EmojiPackEdit)
		emojiPackPrivata.POST("/delete", handler.EmojiPackDelete)
		emojiPackPrivata.GET("/detail", handler.EmojiPackDetail)
		emojiPackPrivata.GET("/listByUser", handler.EmojiPackListByUser)
		emojiPackPrivata.POST("/emojiPackAddEmoji", handler.EmojiPackAddEmoji)
		emojiPackPrivata.POST("/emojiPackRemoveEmoji", handler.EmojiPackRemoveEmoji)
		emojiPackPrivata.POST("/emojiPackCollection", handler.EmojiPackCollection)

	}
	emojiPackPublic := emojiPackGroup.Group("")
	{
		emojiPackPublic.GET("/list", handler.EmojiPackList)
		emojiPackPublic.GET("/emojiPackGetEmojis", handler.EmojiPackGetEmojis)
	}

	return r
}
