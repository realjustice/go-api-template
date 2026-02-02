package logger

import (
	"go-api-template/pkg/config"

	"go.uber.org/zap"
)

// InitLogger 从配置初始化日志
func InitLogger(cfg *config.Config) (*zap.Logger, error) {
	loggerConfig := &Config{
		Level:      cfg.Logger.Level,
		Filename:   cfg.Logger.Filename,
		MaxSize:    cfg.Logger.MaxSize,
		MaxBackups: cfg.Logger.MaxBackups,
		MaxAge:     cfg.Logger.MaxAge,
		Compress:   cfg.Logger.Compress,
		Console:    cfg.Logger.Console,
	}

	return NewLogger(loggerConfig)
}
