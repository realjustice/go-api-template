//go:build wireinject
// +build wireinject

package main

import (
	"go-api-template/internal/controller"
	"go-api-template/internal/middleware"
	"go-api-template/internal/repository"
	"go-api-template/internal/service"
	"go-api-template/pkg/config"
	"go-api-template/pkg/database"
	"go-api-template/pkg/logger"
	"go-api-template/pkg/web"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"
)

// InitializeApp 初始化应用
func InitializeApp(configPath string) (*gin.Engine, func(), error) {
	wire.Build(
		// 配置
		config.LoadConfig,

		// 日志
		logger.InitLogger,

		// 数据库
		database.NewMySQLDB,

		// Repository - Demo 数据访问层
		repository.NewDemoRepository,

		// Service - Demo 业务逻辑层
		service.NewDemoService,

		// Controller - Demo 控制器
		controller.NewDemoController,

		// Middleware - 中间件
		middleware.NewMiddleware,

		// Router - 路由配置和清理函数
		provideRouterAndCleanup,
	)
	return nil, nil, nil
}

// provideRouterAndCleanup 配置路由并提供清理函数
func provideRouterAndCleanup(
	cfg *config.Config,
	demoCtrl *controller.DemoController,
	mw *middleware.Middleware,
	_ *zap.Logger, // 确保 logger 被初始化
) (*gin.Engine, func()) {
	router := provideRouter(cfg, demoCtrl, mw)
	cleanup := func() {
		logger.Close()
	}
	return router, cleanup
}

// provideRouter 配置路由
func provideRouter(
	cfg *config.Config,
	demoCtrl *controller.DemoController,
	mw *middleware.Middleware,
) *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(web.ToGinHandler(mw.RequestID.Handle()))

	// 处理 404 错误
	r.NoRoute(web.ToGinHandler(web.NotFoundHandler()))

	// 处理 405 错误
	r.NoMethod(web.ToGinHandler(web.MethodNotAllowedHandler()))

	// 健康检查（无需鉴权）
	r.GET("/health", web.ToGinHandler(web.HealthHandler()))

	// API v1 路由组
	api := r.Group("/api/v1")
	{
		// Demo CRUD 示例接口
		demos := api.Group("/demos")
		{
			demos.GET("", web.ToGinHandler(demoCtrl.GetAll))       // 获取所有 Demo
			demos.GET("/:id", web.ToGinHandler(demoCtrl.GetByID))  // 获取单个 Demo
			demos.POST("", web.ToGinHandler(demoCtrl.Create))      // 创建 Demo
			demos.PUT("/:id", web.ToGinHandler(demoCtrl.Update))   // 更新 Demo
			demos.DELETE("/:id", web.ToGinHandler(demoCtrl.Delete)) // 删除 Demo
		}
	}

	return r
}
