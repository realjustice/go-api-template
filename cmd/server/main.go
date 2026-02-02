package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-api-template/pkg/config"
	"go-api-template/pkg/logger"
)

func main() {
	// 解析命令行参数
	configPath := flag.String("config", "config/config.yaml", "配置文件路径")
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}

	// 初始化日志
	_, err = logger.InitLogger(cfg)
	if err != nil {
		log.Fatalf("❌ 初始化日志失败: %v", err)
	}
	defer logger.Close()

	logger.Info("🚀 应用启动中...")

	// 初始化应用（通过 Wire 依赖注入）
	router, cleanup, err := InitializeApp(*configPath)
	if err != nil {
		logger.Fatalf("❌ 初始化应用失败: %v", err)
	}
	defer cleanup() // 确保在退出时清理资源

	// 服务器端口
	port := fmt.Sprintf(":%d", cfg.Server.Port)

	// 打印启动信息
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("  Go API Template - 服务已启动")
	fmt.Println("========================================")
	fmt.Printf("🌐 服务地址: http://localhost%s\n", port)
	fmt.Printf("📚 API 文档:\n")
	fmt.Printf("   - 健康检查:    GET  http://localhost%s/health\n", port)
	fmt.Printf("   - Demo 列表:   GET  http://localhost%s/api/v1/demos\n", port)
	fmt.Printf("   - Demo 详情:   GET  http://localhost%s/api/v1/demos/:id\n", port)
	fmt.Printf("   - 创建 Demo:   POST http://localhost%s/api/v1/demos\n", port)
	fmt.Printf("   - 更新 Demo:   PUT  http://localhost%s/api/v1/demos/:id\n", port)
	fmt.Printf("   - 删除 Demo:   DEL  http://localhost%s/api/v1/demos/:id\n", port)
	fmt.Println("========================================")
	fmt.Printf("💡 使用 Ctrl+C 停止服务\n")
	fmt.Println()

	logger.Infof("服务器启动在端口 %s", port)

	// 启动服务器（在 goroutine 中）
	go func() {
		if err := router.Run(port); err != nil {
			logger.Fatalf("❌ 服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号（优雅关闭）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("⏳ 正在关闭服务器...")
	fmt.Println()
	fmt.Println("✅ 服务器已关闭")
}
