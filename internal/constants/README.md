# Constants 层说明

## 📖 包的作用

Constants 包集中管理项目中的所有常量，避免硬编码，提高代码可维护性。

## 🎯 职责范围

### ✅ Constants 应该包含什么

- HTTP Header 常量
- Context Key 常量
- API 响应消息常量
- 业务状态码常量
- 错误消息常量
- 配置默认值常量

### ❌ Constants 不应该包含什么

- 不包含业务逻辑
- 不包含可变的值
- 不包含配置（配置放 config.yaml）

## 📝 现有常量文件

### 1. `context.go` - Context Key 常量

存储在 Context 中的键名：

```go
package constants

const (
    CtxKeyRequestID = "request_id"  // 请求ID
    CtxKeyUserID    = "user_id"     // 用户ID
    CtxKeyUsername  = "username"    // 用户名
)
```

**使用示例：**

```go
// 存储
ctx.Set(constants.CtxKeyUserID, userID)

// 获取
userID := ctx.GetString(constants.CtxKeyUserID)
```

### 2. `header.go` - HTTP Header 常量

HTTP 请求/响应头的键名：

```go
package constants

const (
    HeaderRequestID     = "X-Request-ID"     // 请求ID
    HeaderAuthorization = "Authorization"    // 授权Token
    HeaderContentType   = "Content-Type"     // 内容类型
)
```

**使用示例：**

```go
// 获取 Header
token := ctx.GetHeader(constants.HeaderAuthorization)

// 设置 Header
ctx.Header(constants.HeaderRequestID, requestID)
```

### 3. `message.go` - API 响应消息常量

统一的 API 响应消息：

```go
package constants

const (
    // 通用消息
    MsgSuccess = "success"
    MsgFailed  = "failed"
    
    // 错误消息
    MsgInterfaceNotFound  = "接口不存在"
    MsgBadRequest         = "请求参数错误"
    MsgUnauthorized       = "未授权"
    MsgNotFound           = "资源不存在"
    MsgInternalError      = "服务器内部错误"
)
```

**使用示例：**

```go
// Controller 中使用
web.NotFound(ctx, constants.MsgNotFound)
web.BadRequest(ctx, constants.MsgBadRequest)
```

### 4. `log.go` - 日志字段常量

日志中使用的字段名：

```go
package constants

const (
    LogFieldRequestID = "request_id"
    LogFieldUserID    = "user_id"
    LogFieldAction    = "action"
    LogFieldError     = "error"
)
```

**使用示例：**

```go
logger.Info("user login",
    logger.String(constants.LogFieldUserID, userID),
    logger.String(constants.LogFieldAction, "login"),
)
```

## 🏗️ 添加新的常量

### 1. 业务状态常量

```go
// status.go
package constants

// 用户状态
const (
    UserStatusActive   = 1  // 激活
    UserStatusInactive = 0  // 未激活
    UserStatusBanned   = -1 // 禁用
)

// 订单状态
const (
    OrderStatusPending   = 1  // 待支付
    OrderStatusPaid      = 2  // 已支付
    OrderStatusShipped   = 3  // 已发货
    OrderStatusCompleted = 4  // 已完成
    OrderStatusCancelled = 0  // 已取消
)
```

### 2. 业务类型常量

```go
// type.go
package constants

// 用户类型
const (
    UserTypeNormal = 1  // 普通用户
    UserTypeVIP    = 2  // VIP用户
    UserTypeAdmin  = 9  // 管理员
)

// 支付方式
const (
    PaymentTypeAlipay = "alipay"
    PaymentTypeWechat = "wechat"
    PaymentTypeCard   = "card"
)
```

### 3. 错误码常量

```go
// code.go
package constants

// HTTP 状态码
const (
    CodeSuccess         = 200
    CodeBadRequest      = 400
    CodeUnauthorized    = 401
    CodeForbidden       = 403
    CodeNotFound        = 404
    CodeInternalError   = 500
)

// 业务错误码
const (
    BizCodeSuccess          = 0     // 成功
    BizCodeParamError       = 1001  // 参数错误
    BizCodeUserNotFound     = 2001  // 用户不存在
    BizCodePasswordWrong    = 2002  // 密码错误
    BizCodeOrderNotFound    = 3001  // 订单不存在
)
```

### 4. 配置默认值常量

```go
// default.go
package constants

// 默认配置值
const (
    DefaultPageSize     = 20    // 默认分页大小
    DefaultCacheTTL     = 300   // 默认缓存时间（秒）
    DefaultTokenExpire  = 7200  // Token 过期时间（秒）
    MaxUploadSize       = 10    // 最大上传文件大小（MB）
)
```

## 💡 最佳实践

### 1. 分文件管理

```go
constants/
├── context.go    // Context 相关
├── header.go     // HTTP Header 相关
├── message.go    // 消息常量
├── status.go     // 状态常量
├── code.go       // 错误码
└── README.md     // 本文档
```

### 2. 使用有意义的命名

```go
// ✅ 清晰明了
const UserStatusActive = 1

// ❌ 含义不清
const Status1 = 1
```

### 3. 分组定义

```go
// ✅ 使用 const 块分组
const (
    StatusActive   = 1
    StatusInactive = 0
)

// ❌ 分散定义
const StatusActive = 1
const StatusInactive = 0
```

### 4. 添加注释

```go
const (
    UserStatusActive   = 1  // 激活
    UserStatusInactive = 0  // 未激活
    UserStatusBanned   = -1 // 禁用
)
```

### 5. 使用类型安全的常量

```go
// 定义类型
type UserStatus int

// 定义常量
const (
    UserStatusActive   UserStatus = 1
    UserStatusInactive UserStatus = 0
)

// 使用
func UpdateStatus(status UserStatus) {
    // 类型安全，只能传入 UserStatus 类型
}
```

## 🎨 代码组织

```go
package constants

// 1. 导入（如果需要）
import (
    "time"
)

// 2. 类型定义（如果需要）
type UserStatus int

// 3. 常量定义（分组）
const (
    // Context Key
    CtxKeyUserID = "user_id"
)

const (
    // 用户状态
    UserStatusActive UserStatus = 1
)

const (
    // 默认值
    DefaultPageSize = 20
)
```

## 📋 命名规范

### 前缀命名

```go
// Header 常量
HeaderRequestID
HeaderAuthorization

// Context Key 常量
CtxKeyUserID
CtxKeyRequestID

// 消息常量
MsgSuccess
MsgNotFound

// 日志字段常量
LogFieldUserID
LogFieldAction

// 状态常量
UserStatusActive
OrderStatusPending

// 错误码常量
CodeSuccess
BizCodeParamError
```

### 使用大写和下划线

```go
// ✅ 导出常量：大写开头
const UserStatusActive = 1

// ✅ 私有常量：小写开头
const defaultTimeout = 30

// 如果是缩写，全部大写
const HeaderHTTPSOnly = "HTTPS-Only"
const CtxKeyRequestID = "request_id"
```

## 🔄 与其他层的关系

```
┌─────────────────────────────────────┐
│  所有层都可以使用 Constants          │
├─────────────────────────────────────┤
│  Controller → constants.MsgSuccess  │
│  Service    → constants.CtxKeyUserID│
│  Middleware → constants.HeaderAuth  │
└─────────────────────────────────────┘
```

## 📚 使用示例

### 在 Controller 中

```go
web.Success(ctx, constants.MsgSuccess)
web.NotFound(ctx, constants.MsgNotFound)
```

### 在 Middleware 中

```go
requestID := ctx.GetHeader(constants.HeaderRequestID)
ctx.Set(constants.CtxKeyUserID, userID)
```

### 在 Service 中

```go
logger.Info("user action",
    logger.String(constants.LogFieldUserID, userID),
    logger.String(constants.LogFieldAction, "login"),
)
```

## 📚 相关文档

- [项目结构说明](../../STRUCTURE.md)

---

**Constants 让代码更规范、更易维护！** 📌
