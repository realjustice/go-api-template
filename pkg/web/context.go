package web

import (
	"go-api-template/internal/constants"

	"github.com/gin-gonic/gin"
)

// Map 替代 gin.H，用于构造 JSON 响应数据
// 隔离 Gin 框架依赖
type Map map[string]interface{}

type Context struct {
	*gin.Context
}

// GetRequestID 获取请求ID
func (c *Context) GetRequestID() string {
	reqID := c.GetString(constants.CtxKeyRequestID)
	return reqID
}
