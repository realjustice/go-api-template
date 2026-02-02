# Service 层说明

## 📖 包的作用

Service 层（业务逻辑层）负责实现核心业务逻辑，处于 Controller 和 Repository 之间。

## 🎯 职责范围

### ✅ Service 层应该做什么

- 实现业务规则和逻辑
- 业务数据验证
- 调用 Repository 层进行数据操作
- 处理事务（如需要）
- 业务日志记录
- 错误包装和处理
- 数据转换和聚合

### ❌ Service 层不应该做什么

- 不处理 HTTP 请求/响应
- 不直接操作数据库（通过 Repository）
- 不包含 Gin 相关代码
- 不做参数绑定（Controller 做）

## 📝 示例代码

参考 `demo_service.go`，这是一个标准的 Service 实现。

### 基本结构

```go
package service

import (
    "context"
    
    "go-api-template/internal/model"
    "go-api-template/internal/repository"
    "go-api-template/pkg/errors"
    "go-api-template/pkg/logger"
)

// DemoService Demo 业务逻辑层
type DemoService struct {
    demoRepo *repository.DemoRepository
}

// NewDemoService 创建 Service（依赖注入）
func NewDemoService(demoRepo *repository.DemoRepository) *DemoService {
    return &DemoService{
        demoRepo: demoRepo,
    }
}

// GetByID 根据 ID 获取
func (s *DemoService) GetByID(ctx context.Context, id uint) (*model.Demo, error) {
    // 1. 调用 Repository
    demo, err := s.demoRepo.FindByID(ctx, id)
    if err != nil {
        // 2. 记录日志
        logger.Error("get demo by id failed",
            logger.Uint("id", id),
            logger.Err(err),
        )
        return nil, err
    }
    
    // 3. 返回结果
    return demo, nil
}

// Create 创建
func (s *DemoService) Create(ctx context.Context, demo *model.Demo) error {
    // 1. 业务逻辑校验
    if demo.Title == "" {
        return errors.New("title cannot be empty")
    }
    
    // 2. 调用 Repository
    err := s.demoRepo.Create(ctx, demo)
    if err != nil {
        logger.Error("create demo failed",
            logger.String("title", demo.Title),
            logger.Err(err),
        )
        return err
    }
    
    // 3. 记录成功日志
    logger.Info("demo created successfully",
        logger.Uint("id", demo.ID),
        logger.String("title", demo.Title),
    )
    
    return nil
}
```

## 🏗️ 开发新的 Service

### 1. 文件命名

```
{模块名}_service.go

示例：
user_service.go
order_service.go
product_service.go
```

### 2. 类型定义

```go
type UserService struct {
    userRepo    *repository.UserRepository    // Repository 依赖
    orderRepo   *repository.OrderRepository   // 可以依赖多个 Repository
    cacheClient cache.CacheInterface          // 其他依赖
}
```

### 3. 构造函数

```go
// New{模块名}Service - 依赖通过参数注入
func NewUserService(
    userRepo *repository.UserRepository,
    orderRepo *repository.OrderRepository,
) *UserService {
    return &UserService{
        userRepo:  userRepo,
        orderRepo: orderRepo,
    }
}
```

### 4. 业务方法

```go
// 方法签名规范
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*model.User, error)
func (s *UserService) CreateUser(ctx context.Context, user *model.User) error
func (s *UserService) UpdateUser(ctx context.Context, id uint, user *model.User) error
func (s *UserService) DeleteUser(ctx context.Context, id uint) error
```

## 💡 最佳实践

### 1. 使用 Context

```go
// ✅ 正确：始终传递 context
func (s *UserService) GetUser(ctx context.Context, id uint) (*model.User, error) {
    return s.userRepo.FindByID(ctx, id)
}

// ❌ 错误：不传递 context
func (s *UserService) GetUser(id uint) (*model.User, error) {
    return s.userRepo.FindByID(id)
}
```

### 2. 业务逻辑校验

```go
func (s *UserService) Create(ctx context.Context, user *model.User) error {
    // 业务规则验证
    if user.Email == "" {
        return errors.New("email is required")
    }
    
    // 检查邮箱是否已存在
    existing, _ := s.userRepo.FindByEmail(ctx, user.Email)
    if existing != nil {
        return errors.New("email already exists")
    }
    
    // 执行创建
    return s.userRepo.Create(ctx, user)
}
```

### 3. 错误处理

```go
func (s *UserService) GetByID(ctx context.Context, id uint) (*model.User, error) {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        // 包装错误，添加上下文
        return nil, errors.Wrapf(err, "failed to get user, id: %d", id)
    }
    return user, nil
}
```

### 4. 日志记录

```go
// ✅ 记录关键业务操作
logger.Info("user created",
    logger.Uint("user_id", user.ID),
    logger.String("email", user.Email),
)

// ✅ 记录错误
logger.Error("create user failed",
    logger.String("email", user.Email),
    logger.Err(err),
)

// ❌ 不要记录过多细节日志（影响性能）
```

### 5. 事务处理

```go
func (s *UserService) CreateUserWithOrder(ctx context.Context, user *model.User, order *model.Order) error {
    // 开启事务
    return s.db.Transaction(func(tx *gorm.DB) error {
        // 创建用户
        if err := s.userRepo.CreateWithTx(ctx, tx, user); err != nil {
            return err
        }
        
        // 创建订单
        order.UserID = user.ID
        if err := s.orderRepo.CreateWithTx(ctx, tx, order); err != nil {
            return err
        }
        
        return nil
    })
}
```

### 6. 可复用性

```go
// ✅ Service 方法应该可以被多个 Controller 复用
// Controller A
func (c *UserController) GetUser(ctx *web.Context) {
    user, err := c.userService.GetUserByID(ctx.Request.Context(), id)
    // ...
}

// Controller B
func (c *AdminController) GetUser(ctx *web.Context) {
    user, err := c.userService.GetUserByID(ctx.Request.Context(), id)
    // ...
}
```

## 🔄 与其他层的关系

```
Controller (HTTP 处理)
    ↓ 调用
Service (业务逻辑) ← 你在这里
    ↓ 调用
Repository (数据访问)
    ↓ 操作
Database
```

### Controller → Service

```go
// Controller 调用 Service
user, err := c.userService.GetUserByID(ctx.Request.Context(), id)
```

### Service → Repository

```go
// Service 调用 Repository
user, err := s.userRepo.FindByID(ctx, id)
```

## 📋 方法命名规范

### 查询操作

```go
GetByID(ctx context.Context, id uint) (*model.User, error)
GetAll(ctx context.Context) ([]*model.User, error)
GetByEmail(ctx context.Context, email string) (*model.User, error)
List(ctx context.Context, page, size int) ([]*model.User, int64, error)
```

### 创建操作

```go
Create(ctx context.Context, user *model.User) error
CreateBatch(ctx context.Context, users []*model.User) error
```

### 更新操作

```go
Update(ctx context.Context, id uint, user *model.User) error
UpdateStatus(ctx context.Context, id uint, status int) error
```

### 删除操作

```go
Delete(ctx context.Context, id uint) error
DeleteBatch(ctx context.Context, ids []uint) error
SoftDelete(ctx context.Context, id uint) error  // 软删除
```

## 🎨 代码组织

```go
package service

import (
    // 1. 标准库
    "context"
    "time"
    
    // 2. 项目内部包
    "go-api-template/internal/model"
    "go-api-template/internal/repository"
    
    // 3. pkg 包
    "go-api-template/pkg/errors"
    "go-api-template/pkg/logger"
    
    // 4. 第三方库
    "gorm.io/gorm"
)

// 1. 类型定义
type UserService struct {
    userRepo *repository.UserRepository
}

// 2. 构造函数
func NewUserService(userRepo *repository.UserRepository) *UserService {
    return &UserService{userRepo: userRepo}
}

// 3. 查询方法
func (s *UserService) GetByID(...) {}
func (s *UserService) GetAll(...) {}

// 4. 创建方法
func (s *UserService) Create(...) {}

// 5. 更新方法
func (s *UserService) Update(...) {}

// 6. 删除方法
func (s *UserService) Delete(...) {}

// 7. 其他业务方法
func (s *UserService) ValidateUser(...) {}
```

## 🔗 依赖注入

在 `cmd/server/wire.go` 中注册：

```go
wire.Build(
    // ...
    repository.NewUserRepository,
    service.NewUserService,     // 添加这里
    // ...
)
```

## 📚 相关文档

- [Controller 层说明](../controller/README.md)
- [Repository 层说明](../repository/README.md)
- [项目结构说明](../../STRUCTURE.md)

---

**Service 层是业务逻辑的核心！** 💼
