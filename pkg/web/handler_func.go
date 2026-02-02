package web

import "github.com/gin-gonic/gin"

// HandlerFunc 定义 HTTP 请求处理函数类型
// 隔离 Gin 框架依赖，业务代码使用 web.HandlerFunc 而不是 gin.HandlerFunc
type HandlerFunc func(*Context)

// ToGinHandler 将 web.HandlerFunc 转换为 gin.HandlerFunc
// 用于在路由注册时进行适配
func ToGinHandler(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 将 gin.Context 包装为 web.Context
		ctx := &Context{Context: c}
		handler(ctx)
	}
}

// FromGinHandler 将 gin.HandlerFunc 转换为 web.HandlerFunc
// 用于兼容现有的 Gin 中间件
func FromGinHandler(ginHandler gin.HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		ginHandler(ctx.Context)
	}
}

// ToGinHandlers 批量转换 web.HandlerFunc 为 gin.HandlerFunc
// 用于路由组批量注册中间件
func ToGinHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		ginHandlers[i] = ToGinHandler(handler)
	}
	return ginHandlers
}
