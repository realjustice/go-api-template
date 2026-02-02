# Middleware 层说明

## 📖 包的作用

Middleware 层提供 HTTP 请求的拦截和预处理功能，在请求到达 Controller 之前或响应返回客户端之前进行统一处理。

## 🎯 职责范围

### ✅ Middleware 应该做什么

- 请求预处理（RequestID、日志等）
- 身份认证（Token 验证）
- 权限验证（角色检查）
- 请求签名验证（CheckSum）
- 限流和熔断
- 跨域处理（CORS）
- 请求/响应日志
- 统一错误处理

### ❌ Middleware 不应该做什么

- 不包含业务逻辑
- 不直接操作数据库（通过 Service）
- 保持简洁和高效

## 📝 示例代码

参考 `request_id.go`，这是一个标准的中间件实现。

### 基本结构

```go
package middleware

import (
    "go-api-template/internal/constants"
    "go-api-template/pkg/web"
    
    "github.com/google/uuid"
)

// RequestIDMiddleware RequestID 中间件
type RequestIDMiddleware struct{}

// NewRequestIDMiddleware 创建中间件
func NewRequestIDMiddleware() *RequestIDMiddleware {
    return &RequestIDMiddleware{}
}

// Handle 处理函数
func (m *RequestIDMiddleware) Handle() web.HandlerFunc {
    return func(ctx *web.Context) {
        // 1. 预处理：尝试从 Header 获取 RequestID
        requestID := ctx.GetHeader(constants.HeaderRequestID)
        
        // 2. 如果没有，生成新的
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        // 3. 存入 Context
        ctx.Set(constants.CtxKeyRequestID, requestID)
        
        // 4. 写入响应头
        ctx.Header(constants.HeaderRequestID, requestID)
        
        // 5. 继续处理请求
        ctx.Next()
    }
}
```

## 🏗️ 开发新的 Middleware

### 1. 文件命名

```
{功能名}.go

示例：
auth.go           - 认证中间件
permission.go     - 权限中间件
rate_limit.go     - 限流中间件
cors.go           - 跨域中间件
```

### 2. 认证中间件示例

```go
package middleware

import (
    "strings"
    
    "go-api-template/internal/constants"
    "go-api-template/pkg/errors"
    "go-api-template/pkg/web"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
    // 可以注入 Service 依赖
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware() *AuthMiddleware {
    return &AuthMiddleware{}
}

// Handle 处理函数
func (m *AuthMiddleware) Handle() web.HandlerFunc {
    return func(ctx *web.Context) {
        // 1. 获取 Token
        token := ctx.GetHeader(constants.HeaderAuthorization)
        if token == "" {
            web.Unauthorized(ctx, "missing token")
            ctx.Abort()  // 中断请求
            return
        }
        
        // 2. 验证 Token（示例）
        token = strings.TrimPrefix(token, "Bearer ")
        userID, err := m.validateToken(token)
        if err != nil {
            web.Unauthorized(ctx, "invalid token")
            ctx.Abort()
            return
        }
        
        // 3. 存入 Context
        ctx.Set(constants.CtxKeyUserID, userID)
        
        // 4. 继续处理
        ctx.Next()
    }
}

func (m *AuthMiddleware) validateToken(token string) (string, error) {
    // Token 验证逻辑
    // ...
    return "user_id", nil
}
```

### 3. 日志中间件示例

```go
package middleware

import (
    "time"
    
    "go-api-template/pkg/logger"
    "go-api-template/pkg/web"
)

// LoggerMiddleware 日志中间件
type LoggerMiddleware struct{}

func NewLoggerMiddleware() *LoggerMiddleware {
    return &LoggerMiddleware{}
}

func (m *LoggerMiddleware) Handle() web.HandlerFunc {
    return func(ctx *web.Context) {
        // 记录开始时间
        start := time.Now()
        path := ctx.Request.URL.Path
        method := ctx.Request.Method
        
        // 处理请求
        ctx.Next()
        
        // 记录日志
        elapsed := time.Since(start)
        logger.Info("HTTP Request",
            logger.String("method", method),
            logger.String("path", path),
            logger.Int("status", ctx.Writer.Status()),
            logger.Duration("elapsed", elapsed),
        )
    }
}
```

### 4. CORS 中间件示例

```go
package middleware

import "go-api-template/pkg/web"

// CORSMiddleware 跨域中间件
type CORSMiddleware struct{}

func NewCORSMiddleware() *CORSMiddleware {
    return &CORSMiddleware{}
}

func (m *CORSMiddleware) Handle() web.HandlerFunc {
    return func(ctx *web.Context) {
        ctx.Header("Access-Control-Allow-Origin", "*")
        ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
        
        // OPTIONS 请求直接返回
        if ctx.Request.Method == "OPTIONS" {
            ctx.AbortWithStatus(204)
            return
        }
        
        ctx.Next()
    }
}
```

## 💡 最佳实践

### 1. 中间件的执行顺序

```go
// wire.go 中的注册顺序很重要
r.Use(gin.Recovery())           // 1. 最外层：Recovery（捕获 panic）
r.Use(LoggerMiddleware())       // 2. 日志记录
r.Use(CORSMiddleware())         // 3. CORS
r.Use(RequestIDMiddleware())    // 4. RequestID
r.Use(AuthMiddleware())         // 5. 认证（最内层）
```

### 2. 使用 Next() 和 Abort()

```go
// ✅ 继续处理
ctx.Next()

// ✅ 中断请求（不执行后续中间件和 Handler）
ctx.Abort()
ctx.AbortWithStatus(401)
```

### 3. 存储和获取数据

```go
// 存储数据
ctx.Set(constants.CtxKeyUserID, userID)

// 在 Controller 中获取
userID := ctx.GetString(constants.CtxKeyUserID)
```

### 4. 可选中间件

```go
// 某些路由需要，某些不需要
api := r.Group("/api/v1")
{
    // 公开接口（无需认证）
    api.GET("/public", handler)
    
    // 需要认证的接口
    auth := api.Group("")
    auth.Use(web.ToGinHandler(authMiddleware.Handle()))
    {
        auth.GET("/users", handler)
        auth.POST("/orders", handler)
    }
}
```

### 5. 中间件依赖注入

```go
// AuthMiddleware 可能需要 Service 依赖
type AuthMiddleware struct {
    tokenService *service.TokenService  // 注入 Service
}

func NewAuthMiddleware(tokenService *service.TokenService) *AuthMiddleware {
    return &AuthMiddleware{
        tokenService: tokenService,
    }
}

// 在 wire.go 中
wire.Build(
    service.NewTokenService,
    middleware.NewAuthMiddleware,  // Wire 会自动注入依赖
)
```

## 📋 常见中间件

### 1. RequestID - 请求追踪

- 为每个请求生成唯一 ID
- 便于日志追踪和问题排查

### 2. Auth - 身份认证

- Token 验证
- 用户身份识别
- 将用户信息存入 Context

### 3. Permission - 权限验证

- 检查用户权限
- 基于角色的访问控制（RBAC）

### 4. RateLimit - 限流

- 防止 API 滥用
- 保护服务器资源

### 5. CORS - 跨域

- 处理跨域请求
- 配置允许的域名和方法

### 6. Logger - 请求日志

- 记录请求信息
- 记录响应时间
- 便于监控和调试

## 🎨 中间件集合

在 `middleware.go` 中统一管理：

```go
package middleware

// Middleware 中间件集合
type Middleware struct {
    RequestID  *RequestIDMiddleware
    Auth       *AuthMiddleware
    Permission *PermissionMiddleware
    RateLimit  *RateLimitMiddleware
    CORS       *CORSMiddleware
}

// NewMiddleware 创建中间件集合
func NewMiddleware(
    tokenService *service.TokenService,
) *Middleware {
    return &Middleware{
        RequestID:  NewRequestIDMiddleware(),
        Auth:       NewAuthMiddleware(tokenService),
        Permission: NewPermissionMiddleware(),
        RateLimit:  NewRateLimitMiddleware(),
        CORS:       NewCORSMiddleware(),
    }
}
```

## 🔗 在路由中使用

```go
// wire.go
func provideRouter(mw *middleware.Middleware) *gin.Engine {
    r := gin.New()
    
    // 全局中间件
    r.Use(gin.Recovery())
    r.Use(web.ToGinHandler(mw.RequestID.Handle()))
    r.Use(web.ToGinHandler(mw.CORS.Handle()))
    
    // 公开 API
    r.GET("/public", handler)
    
    // 需要认证的 API
    api := r.Group("/api/v1")
    api.Use(web.ToGinHandler(mw.Auth.Handle()))
    {
        api.GET("/users", handler)
    }
    
    return r
}
```

## 🔗 依赖注入

在 `cmd/server/wire.go` 中注册：

```go
wire.Build(
    // ...
    middleware.NewRequestIDMiddleware,
    middleware.NewAuthMiddleware,
    middleware.NewMiddleware,  // 中间件集合
    // ...
)
```

## 📚 相关文档

- [Controller 层说明](../controller/README.md)
- [项目结构说明](../../STRUCTURE.md)

---

**中间件是请求处理的第一道关卡！** 🛡️
