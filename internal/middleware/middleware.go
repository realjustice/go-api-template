package middleware

import (
	"go-api-template/pkg/config"
)

// Middleware 中间件集合
type Middleware struct {
	RequestID *RequestIDMiddleware
	CORS      *CORSMiddleware
}

// NewMiddleware 创建中间件集合
func NewMiddleware(cfg *config.Config) *Middleware {
	// 根据配置创建 CORS 中间件
	var corsMiddleware *CORSMiddleware
	if cfg.CORS.Enabled {
		corsMiddleware = NewCORSMiddleware(&CORSConfig{
			AllowOrigins: cfg.CORS.AllowOrigins,
			AllowMethods: cfg.CORS.AllowMethods,
			AllowHeaders: cfg.CORS.AllowHeaders,
		})
	} else {
		// CORS 未启用时使用默认配置
		corsMiddleware = NewDefaultCORSMiddleware()
	}

	return &Middleware{
		RequestID: NewRequestIDMiddleware(),
		CORS:      corsMiddleware,
	}
}
