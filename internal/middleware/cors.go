package middleware

import (
	"go-api-template/pkg/web"
)

// CORSMiddleware CORS 跨域中间件
type CORSMiddleware struct {
	allowOrigins []string
	allowMethods []string
	allowHeaders []string
}

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowOrigins []string // 允许的来源，如：["http://localhost:3000", "https://example.com"]
	AllowMethods []string // 允许的方法，如：["GET", "POST", "PUT", "DELETE"]
	AllowHeaders []string // 允许的请求头，如：["Content-Type", "Authorization"]
}

// NewCORSMiddleware 创建 CORS 中间件
func NewCORSMiddleware(config *CORSConfig) *CORSMiddleware {
	// 设置默认值
	if config == nil {
		config = &CORSConfig{}
	}

	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = []string{"*"} // 默认允许所有来源
	}

	if len(config.AllowMethods) == 0 {
		config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	}

	if len(config.AllowHeaders) == 0 {
		config.AllowHeaders = []string{"Content-Type", "Authorization", "X-Request-ID"}
	}

	return &CORSMiddleware{
		allowOrigins: config.AllowOrigins,
		allowMethods: config.AllowMethods,
		allowHeaders: config.AllowHeaders,
	}
}

// NewDefaultCORSMiddleware 创建默认配置的 CORS 中间件
func NewDefaultCORSMiddleware() *CORSMiddleware {
	return NewCORSMiddleware(nil)
}

// Handle CORS 处理函数
func (m *CORSMiddleware) Handle() web.HandlerFunc {
	return func(ctx *web.Context) {
		// 获取请求来源
		origin := ctx.GetHeader("Origin")

		// 检查来源是否允许
		if m.isOriginAllowed(origin) {
			// 设置 CORS 响应头
			ctx.Header("Access-Control-Allow-Origin", origin)
		} else if len(m.allowOrigins) == 1 && m.allowOrigins[0] == "*" {
			// 允许所有来源
			ctx.Header("Access-Control-Allow-Origin", "*")
		}

		// 设置其他 CORS 响应头
		ctx.Header("Access-Control-Allow-Methods", m.joinStrings(m.allowMethods))
		ctx.Header("Access-Control-Allow-Headers", m.joinStrings(m.allowHeaders))
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Max-Age", "86400") // 预检请求缓存 24 小时

		// OPTIONS 请求直接返回（预检请求）
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		// 继续处理请求
		ctx.Next()
	}
}

// isOriginAllowed 检查来源是否允许
func (m *CORSMiddleware) isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	for _, allowed := range m.allowOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}

	return false
}

// joinStrings 连接字符串数组
func (m *CORSMiddleware) joinStrings(arr []string) string {
	if len(arr) == 0 {
		return ""
	}

	result := arr[0]
	for i := 1; i < len(arr); i++ {
		result += ", " + arr[i]
	}

	return result
}
