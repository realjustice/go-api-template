package middleware

import (
	"go-api-template/internal/constants"
	"go-api-template/pkg/web"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware RequestID 中间件
type RequestIDMiddleware struct{}

// NewRequestIDMiddleware 创建 RequestID 中间件
func NewRequestIDMiddleware() *RequestIDMiddleware {
	return &RequestIDMiddleware{}
}

// Handle 处理 RequestID
func (m *RequestIDMiddleware) Handle() web.HandlerFunc {
	return func(ctx *web.Context) {
		// 尝试从 Header 获取 RequestID
		requestID := ctx.GetHeader(constants.HeaderRequestID)
		
		// 如果 Header 中没有，则生成新的 UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		// 存入 Context，供后续使用
		ctx.Set(constants.CtxKeyRequestID, requestID)
		
		// 将 RequestID 写入响应头，方便追踪
		ctx.Header(constants.HeaderRequestID, requestID)
		
		ctx.Next()
	}
}

// GetRequestID 从 Context 中获取 RequestID（兼容方法）
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(constants.CtxKeyRequestID); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
