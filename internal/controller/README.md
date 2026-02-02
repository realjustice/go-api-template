# Controller 层说明

## 示例控制器

本目录包含 `demo_controller.go` 作为标准示例，展示了如何实现一个完整的 RESTful API 控制器。

## DemoController 功能

Demo Controller 实现了完整的 CRUD 操作：

- `GetAll()` - 获取所有记录
- `GetByID()` - 根据 ID 获取单条记录
- `Create()` - 创建新记录
- `Update()` - 更新记录
- `Delete()` - 删除记录

## 开发新的 Controller

参考 `demo_controller.go` 创建新的控制器时，请遵循以下规范：

### 1. 命名规范

```go
// 文件名：{模块名}_controller.go
// 类型名：{模块名}Controller
type UserController struct {
    userService *service.UserService
}
```

### 2. 构造函数

```go
func NewUserController(userService *service.UserService) *UserController {
    return &UserController{
        userService: userService,
    }
}
```

### 3. Handler 方法

```go
func (c *UserController) GetByID(ctx *web.Context) {
    // 1. 参数验证
    id := ctx.Param("id")
    
    // 2. 调用 Service
    user, err := c.userService.GetByID(ctx.Request.Context(), id)
    
    // 3. 错误处理
    if err != nil {
        if errors.Is(err, errors.ErrNotFound) {
            web.NotFound(ctx, "用户不存在")
            return
        }
        web.InternalError(ctx, "获取用户失败")
        return
    }
    
    // 4. 返回响应
    web.Success(ctx, user)
}
```

### 4. 请求/响应结构

```go
// 请求结构
type CreateUserRequest struct {
    Name  string `json:"name" binding:"required"`
    Email string `json:"email" binding:"required,email"`
}

// 在 Handler 中使用
var req CreateUserRequest
if err := ctx.ShouldBindJSON(&req); err != nil {
    web.BadRequest(ctx, "参数错误: "+err.Error())
    return
}
```

## 最佳实践

1. **单一职责**: Controller 只负责 HTTP 处理，业务逻辑放在 Service 层
2. **参数验证**: 使用 `binding` 标签进行参数校验
3. **错误处理**: 根据错误类型返回合适的 HTTP 状态码
4. **统一响应**: 使用 `web.Success()` 等方法统一响应格式
5. **上下文传递**: 使用 `ctx.Request.Context()` 传递上下文到下层

## 注册路由

在 `cmd/server/wire.go` 中注册路由：

```go
func provideRouter(
    userCtrl *controller.UserController,
) *gin.Engine {
    r := gin.New()
    
    api := r.Group("/api/v1")
    {
        users := api.Group("/users")
        {
            users.GET("", web.ToGinHandler(userCtrl.GetAll))
            users.GET("/:id", web.ToGinHandler(userCtrl.GetByID))
            users.POST("", web.ToGinHandler(userCtrl.Create))
            users.PUT("/:id", web.ToGinHandler(userCtrl.Update))
            users.DELETE("/:id", web.ToGinHandler(userCtrl.Delete))
        }
    }
    
    return r
}
```

## 依赖注入

在 `cmd/server/wire.go` 的 `InitializeApp` 中添加：

```go
wire.Build(
    // ...
    repository.NewUserRepository,
    service.NewUserService,
    controller.NewUserController,  // 添加新的 Controller
    // ...
)
```
