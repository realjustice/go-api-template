# 项目结构说明

本文档详细介绍 Go API Template 的目录结构和各个包的职责。

## 📁 完整目录树

```
go-api-template/
├── cmd/
│   └── server/              # 应用入口
│       ├── main.go          # 主函数
│       ├── wire.go          # Wire 依赖注入配置
│       └── wire_gen.go      # Wire 生成代码（自动生成）
│
├── internal/                # 私有应用代码
│   ├── controller/          # HTTP 控制器层
│   │   ├── demo_controller.go
│   │   └── README.md
│   │
│   ├── service/             # 业务逻辑层
│   │   └── demo_service.go
│   │
│   ├── repository/          # 数据访问层
│   │   └── demo_repository.go
│   │
│   ├── model/               # 数据模型
│   │   └── demo.go
│   │
│   ├── middleware/          # 中间件
│   │   ├── middleware.go
│   │   └── request_id.go
│   │
│   ├── router/              # 路由配置（可选）
│   │   └── router.go
│   │
│   └── constants/           # 常量定义
│       ├── context.go       # Context Key 常量
│       ├── header.go        # HTTP Header 常量
│       ├── message.go       # API 消息常量
│       └── log.go           # 日志字段常量
│
├── pkg/                     # 公共库（可被外部项目导入）
│   ├── config/              # 配置管理
│   │   └── config.go
│   │
│   ├── database/            # 数据库
│   │   └── mysql.go
│   │
│   ├── redis/               # Redis 客户端
│   │   └── redis.go
│   │
│   ├── cache/               # 缓存门面
│   │   ├── interface.go
│   │   ├── facade.go
│   │   └── factory.go
│   │
│   ├── logger/              # 日志组件
│   │   ├── logger.go
│   │   └── init.go
│   │
│   ├── errors/              # 错误处理
│   │   └── errors.go
│   │
│   ├── web/                 # Web 框架隔离
│   │   ├── context.go
│   │   ├── handler_func.go
│   │   ├── response.go
│   │   └── handlers.go
│   │
│   ├── security/            # 安全工具
│   │   └── checksum.go
│   │
│   └── tools/               # 工具函数
│       └── random.go
│
├── config/                  # 配置文件
│   └── config.yaml
│
├── logs/                    # 日志文件（自动生成）
│   └── app.log
│
├── bin/                     # 编译输出（自动生成）
│   └── server
│
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
├── create-project.sh        # 项目创建脚本
├── README.md                # 项目说明
└── STRUCTURE.md             # 本文档
```

---

## 📖 分层架构说明

### 请求处理流程

```
HTTP Request
     ↓
Controller (HTTP 处理)
     ↓
Service (业务逻辑)
     ↓
Repository (数据访问)
     ↓
Database/Cache
```

---

## 🗂️ 目录详细说明

### 1️⃣ `cmd/server/` - 应用入口

**作用**：应用程序启动入口

**文件说明**：
- `main.go` - 主函数，初始化和启动服务器
- `wire.go` - Wire 依赖注入配置（手动编写）
- `wire_gen.go` - Wire 生成的代码（自动生成，不要手动修改）

**职责**：
- 加载配置文件
- 初始化日志系统
- 通过 Wire 注入依赖
- 启动 HTTP 服务
- 处理优雅关闭

---

### 2️⃣ `internal/` - 私有应用代码

> `internal/` 目录下的代码不能被外部项目导入（Go 语言特性）

#### `controller/` - HTTP 控制器层

**作用**：处理 HTTP 请求和响应

**职责**：
- 接收 HTTP 请求
- 参数验证和绑定
- 调用 Service 层处理业务
- 返回统一格式的响应

**示例**：
```go
func (c *DemoController) GetByID(ctx *web.Context) {
    id := ctx.Param("id")
    demo, err := c.demoService.GetByID(ctx.Request.Context(), id)
    if err != nil {
        web.NotFound(ctx, "demo not found")
        return
    }
    web.Success(ctx, demo)
}
```

**原则**：
- 不包含业务逻辑
- 不直接操作数据库
- 只负责 HTTP 相关处理

---

#### `service/` - 业务逻辑层

**作用**：实现核心业务逻辑

**职责**：
- 业务规则验证
- 业务流程控制
- 事务管理
- 调用 Repository 层
- 错误处理和日志记录

**示例**：
```go
func (s *DemoService) Create(ctx context.Context, demo *model.Demo) error {
    // 业务校验
    if demo.Title == "" {
        return errors.New("title cannot be empty")
    }
    
    // 调用 Repository
    err := s.demoRepo.Create(ctx, demo)
    if err != nil {
        logger.Error("create demo failed", logger.Err(err))
        return err
    }
    
    logger.Info("demo created", logger.Uint("id", demo.ID))
    return nil
}
```

**原则**：
- 可以被多个 Controller 复用
- 不包含 HTTP 相关代码
- 处理业务逻辑和事务

---

#### `repository/` - 数据访问层

**作用**：封装数据库操作

**职责**：
- 数据库 CRUD 操作
- 查询条件封装
- 数据模型转换
- 数据库错误处理

**示例**：
```go
func (r *DemoRepository) FindByID(ctx context.Context, id uint) (*model.Demo, error) {
    var demo model.Demo
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&demo).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.Wrapf(errors.ErrNotFound, "id: %d", id)
        }
        return nil, errors.Wrap(err, "query failed")
    }
    return &demo, nil
}
```

**原则**：
- 只负责数据访问
- 不包含业务逻辑
- 返回数据模型或错误

---

#### `model/` - 数据模型

**作用**：定义数据结构

**职责**：
- 数据库表结构定义
- JSON 序列化标签
- GORM 标签配置
- 表关联定义

**示例**：
```go
type Demo struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Title     string    `json:"title" gorm:"type:varchar(200);not null"`
    Content   string    `json:"content" gorm:"type:text"`
    CreatedAt time.Time `json:"created_at"`
}

func (Demo) TableName() string {
    return "demos"
}
```

**原则**：
- 只包含数据定义
- 不包含业务逻辑

---

#### `middleware/` - 中间件

**作用**：HTTP 请求拦截和处理

**职责**：
- 请求预处理
- 身份认证
- 权限验证
- 日志记录
- 请求追踪

**示例**：
```go
func (m *RequestIDMiddleware) Handle() web.HandlerFunc {
    return func(ctx *web.Context) {
        requestID := ctx.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        ctx.Set("request_id", requestID)
        ctx.Next()
    }
}
```

**原则**：
- 可复用的横切关注点
- 不包含业务逻辑

---

#### `constants/` - 常量定义

**作用**：统一管理常量

**文件说明**：
- `context.go` - Context Key 常量（如：`request_id`, `user_id`）
- `header.go` - HTTP Header 常量（如：`X-Request-ID`）
- `message.go` - API 响应消息常量（如：错误提示）
- `log.go` - 日志字段常量（如：日志字段名）

**职责**：
- 避免硬编码
- 统一管理字符串常量
- 便于维护和修改

---

### 3️⃣ `pkg/` - 公共库

> `pkg/` 目录下的代码可以被外部项目导入

#### `config/` - 配置管理

**作用**：加载和管理应用配置

**功能**：
- 从 YAML 文件加载配置
- 支持多环境配置
- 配置验证和默认值

---

#### `database/` - 数据库

**作用**：数据库连接管理

**功能**：
- MySQL/PostgreSQL 连接
- 连接池配置
- GORM 初始化

---

#### `redis/` - Redis 客户端

**作用**：Redis 连接管理

**功能**：
- Redis 连接
- 连接池配置
- 健康检查

---

#### `cache/` - 缓存门面

**作用**：统一的缓存接口

**功能**：
- 支持多种缓存驱动（Redis/Memory/Chain）
- 统一的 Get/Set/Delete 接口
- 易于切换缓存实现

---

#### `logger/` - 日志组件

**作用**：企业级日志系统

**功能**：
- Zap 高性能日志
- 日志自动切割（Lumberjack）
- 结构化日志
- 多输出（文件 + 控制台）

**使用示例**：
```go
logger.Info("操作成功",
    logger.String("user_id", userID),
    logger.Duration("elapsed", elapsed),
)
```

---

#### `errors/` - 错误处理

**作用**：企业级错误处理

**功能**：
- 完整堆栈跟踪
- 错误链支持
- 错误上下文信息
- 预定义业务错误

**使用示例**：
```go
// 创建错误
err := errors.New("操作失败")

// 包装错误
err = errors.Wrapf(err, "user_id: %s", userID)

// 判断错误类型
if errors.Is(err, errors.ErrNotFound) {
    // 处理
}
```

---

#### `web/` - Web 框架隔离层

**作用**：隔离 Gin 框架依赖

**功能**：
- Context 封装
- HandlerFunc 封装
- 统一响应格式
- 框架无关的业务代码

**优势**：
- 业务代码不直接依赖 Gin
- 易于测试
- 易于迁移到其他框架

**使用示例**：
```go
// 定义 Handler（不依赖 Gin）
func Handler(ctx *web.Context) {
    web.Success(ctx, data)
}

// 注册路由
router.GET("/path", web.ToGinHandler(Handler))
```

---

#### `security/` - 安全工具

**作用**：安全相关工具函数

**功能**：
- 签名验证
- 加密解密
- 哈希计算

---

#### `tools/` - 工具函数

**作用**：通用工具函数

**功能**：
- 随机字符串生成
- 时间处理
- 字符串处理
- 其他通用工具

---

## 🎯 开发流程

### 添加新功能的步骤

#### 1. 定义数据模型

```bash
# 创建 Model
vim internal/model/user.go
```

```go
type User struct {
    ID    uint   `json:"id" gorm:"primaryKey"`
    Name  string `json:"name" gorm:"type:varchar(100)"`
    Email string `json:"email" gorm:"type:varchar(100);uniqueIndex"`
}
```

#### 2. 实现 Repository

```bash
# 创建 Repository
vim internal/repository/user_repository.go
```

```go
type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
    // 实现数据访问逻辑
}
```

#### 3. 实现 Service

```bash
# 创建 Service
vim internal/service/user_service.go
```

```go
type UserService struct {
    userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
    return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
    // 实现业务逻辑
}
```

#### 4. 实现 Controller

```bash
# 创建 Controller
vim internal/controller/user_controller.go
```

```go
type UserController struct {
    userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
    return &UserController{userService: userService}
}

func (c *UserController) GetUser(ctx *web.Context) {
    // 实现 HTTP 处理
}
```

#### 5. 注册依赖和路由

```bash
# 编辑 Wire 配置
vim cmd/server/wire.go
```

在 `InitializeApp` 中添加：
```go
wire.Build(
    // ...
    repository.NewUserRepository,  // 添加 Repository
    service.NewUserService,        // 添加 Service
    controller.NewUserController,  // 添加 Controller
    // ...
)
```

在 `provideRouter` 中添加路由：
```go
users := api.Group("/users")
{
    users.GET("/:id", web.ToGinHandler(userCtrl.GetUser))
}
```

#### 6. 生成代码并运行

```bash
make wire   # 生成 Wire 代码
make run    # 运行项目
```

---

## 📦 包依赖关系

```
┌─────────────────────────────────────────┐
│              cmd/server                  │
│         (应用入口 + Wire配置)             │
└────────────┬────────────────────────────┘
             │ 依赖
             ↓
┌────────────────────────────────────────────┐
│              internal/                      │
│  ┌──────────────────────────────────────┐  │
│  │  controller → service → repository   │  │
│  └──────────────────────────────────────┘  │
│         middleware, constants               │
└────────────┬───────────────────────────────┘
             │ 依赖
             ↓
┌────────────────────────────────────────────┐
│                 pkg/                        │
│  config, database, redis, cache,            │
│  logger, errors, web, security, tools       │
└─────────────────────────────────────────────┘
             │ 依赖
             ↓
┌─────────────────────────────────────────────┐
│            第三方库                          │
│  Gin, GORM, Zap, Wire, Redis, etc.          │
└─────────────────────────────────────────────┘
```

---

## 🔐 包的可见性

### `internal/` 包

- ❌ **不能**被外部项目导入
- ✅ 只能在本项目内部使用
- 用于存放项目特定的业务代码

### `pkg/` 包

- ✅ **可以**被外部项目导入
- 用于存放可复用的公共库
- 需要保持稳定的 API

---

## 💡 最佳实践

### 1. 保持分层清晰

```
✅ Controller 调用 Service
✅ Service 调用 Repository
❌ Controller 不能直接调用 Repository
❌ Repository 不能调用 Service
```

### 2. 单向依赖

```
✅ internal/ 依赖 pkg/
❌ pkg/ 不能依赖 internal/
```

### 3. 错误处理

```go
// Repository 层
return nil, errors.Wrap(err, "query failed")

// Service 层
return nil, errors.Wrapf(err, "user_id: %s", userID)

// Controller 层
if errors.Is(err, errors.ErrNotFound) {
    web.NotFound(ctx, "用户不存在")
}
```

### 4. 日志记录

```go
// Service 层记录业务日志
logger.Info("user created",
    logger.Uint("user_id", user.ID),
    logger.String("name", user.Name),
)

// Controller 层只处理 HTTP 响应
web.Success(ctx, user)
```

### 5. 使用常量

```go
// ❌ 硬编码
ctx.Set("request_id", id)

// ✅ 使用常量
ctx.Set(constants.CtxKeyRequestID, id)
```

---

## 📝 命名规范

### 文件命名

- Controller: `{模块名}_controller.go` (如：`user_controller.go`)
- Service: `{模块名}_service.go` (如：`user_service.go`)
- Repository: `{模块名}_repository.go` (如：`user_repository.go`)
- Model: `{模块名}.go` (如：`user.go`)

### 类型命名

```go
type UserController struct {}   // Controller 后缀
type UserService struct {}      // Service 后缀
type UserRepository struct {}   // Repository 后缀
type User struct {}             // Model 无后缀
```

### 函数命名

```go
// Constructor
func NewUserController() *UserController {}

// Repository 方法
func (r *UserRepository) FindByID() {}
func (r *UserRepository) Create() {}

// Service 方法
func (s *UserService) GetUserByID() {}
func (s *UserService) CreateUser() {}

// Controller 方法
func (c *UserController) GetUser() {}
func (c *UserController) CreateUser() {}
```

---

## 🎨 代码组织建议

### Controller 层

```go
// 1. 类型定义
type UserController struct {
    userService *service.UserService
}

// 2. 构造函数
func NewUserController(...) *UserController {}

// 3. HTTP Handler 方法
func (c *UserController) GetUser(ctx *web.Context) {}
func (c *UserController) CreateUser(ctx *web.Context) {}

// 4. 请求/响应结构（如需要）
type CreateUserRequest struct {}
```

### Service 层

```go
// 1. 类型定义
type UserService struct {
    userRepo *repository.UserRepository
}

// 2. 构造函数
func NewUserService(...) *UserService {}

// 3. 业务方法
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {}
func (s *UserService) CreateUser(ctx context.Context, user *model.User) error {}
```

### Repository 层

```go
// 1. 类型定义
type UserRepository struct {
    db *gorm.DB
}

// 2. 构造函数
func NewUserRepository(db *gorm.DB) *UserRepository {}

// 3. 数据访问方法
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {}
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {}
```

---

## 📚 相关文档

- [README.md](./README.md) - 项目简介和快速开始
- [Controller 开发指南](./internal/controller/README.md) - 详细的 Controller 开发说明

---

**详细了解项目结构，快速上手开发！** 🚀
