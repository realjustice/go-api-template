package web

import (
	"go-api-template/internal/constants"
)

// ========== 常用 Handler 函数 ==========

// HealthHandler 健康检查 Handler
// 返回服务状态，用于负载均衡器和监控
func HealthHandler() HandlerFunc {
	return func(ctx *Context) {
		Success(ctx, Map{
			"status": "ok",
		})
	}
}

// NotFoundHandler 404 错误 Handler
// 返回统一的 JSON 格式 404 响应
func NotFoundHandler() HandlerFunc {
	return func(ctx *Context) {
		NotFound(ctx, constants.MsgInterfaceNotFound)
	}
}

// MethodNotAllowedHandler 405 错误 Handler
// 返回统一的 JSON 格式 405 响应
func MethodNotAllowedHandler() HandlerFunc {
	return func(ctx *Context) {
		Error(ctx, 405, 405, constants.MsgMethodNotAllowed)
	}
}
