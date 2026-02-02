# Repository 层说明

## 📖 包的作用

Repository 层（数据访问层）负责封装所有数据库操作，为 Service 层提供数据访问接口。

## 🎯 职责范围

### ✅ Repository 层应该做什么

- 封装数据库 CRUD 操作
- 构建查询条件
- 处理数据库错误
- 数据模型转换
- 执行原生 SQL（如需要）
- 处理数据库事务

### ❌ Repository 层不应该做什么

- 不包含业务逻辑
- 不做业务数据验证
- 不记录业务日志（只记录数据库操作日志）
- 不处理 HTTP 请求

## 📝 示例代码

参考 `demo_repository.go`，这是一个标准的 Repository 实现。

### 基本结构（使用 BaseRepository）

```go
package repository

import (
    "context"
    
    "go-api-template/internal/model"
    "go-api-template/pkg/database"
    "go-api-template/pkg/errors"
    
    "gorm.io/gorm"
)

// DemoRepository Demo 数据访问层
type DemoRepository struct {
    *database.BaseRepository  // 嵌入 BaseRepository，复用通用方法
    db                       *gorm.DB  // 保留 db 用于复杂查询
}

// NewDemoRepository 创建 Repository（依赖注入）
func NewDemoRepository(db *gorm.DB) *DemoRepository {
    return &DemoRepository{
        BaseRepository: database.NewBaseRepository(db),
        db:             db,
    }
}

// FindByID 根据 ID 查询（使用基类方法）
func (r *DemoRepository) FindByID(ctx context.Context, id uint) (*model.Demo, error) {
    var demo model.Demo
    err := r.BaseRepository.FindByID(ctx, id, &demo)  // 使用基类方法
    if err != nil {
        return nil, errors.Wrapf(err, "demo not found, id: %d", id)
    }
    return &demo, nil
}

// Create 创建（使用基类方法）
func (r *DemoRepository) Create(ctx context.Context, demo *model.Demo) error {
    return r.BaseRepository.Create(ctx, demo)  // 使用基类方法
}

// Search 搜索（复杂查询，直接使用 GORM）
func (r *DemoRepository) Search(ctx context.Context, keyword string) ([]*model.Demo, error) {
    var demos []*model.Demo
    
    // 复杂查询直接用 GORM，保留灵活性
    err := r.db.WithContext(ctx).
        Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
        Where("status = ?", 1).
        Order("created_at DESC").
        Find(&demos).Error
    
    if err != nil {
        return nil, errors.Wrap(err, "search demos failed")
    }
    return demos, nil
}
```

### BaseRepository 优势

✅ **减少重复代码** - 通用 CRUD 操作无需重复实现
✅ **保留灵活性** - 复杂查询仍可直接使用 GORM
✅ **统一错误处理** - 基类统一包装错误
✅ **可选使用** - 通过嵌入方式，不是强制的

## 🏗️ 开发新的 Repository

### 1. 文件命名

```
{模块名}_repository.go

示例：
user_repository.go
order_repository.go
product_repository.go
```

### 2. 类型定义

```go
type UserRepository struct {
    db *gorm.DB  // 必须有数据库连接
}
```

### 3. 构造函数

```go
func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}
```

### 4. 查询方法

```go
// 单条查询
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
    var user model.User
    err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.Wrapf(errors.ErrNotFound, "user not found, id: %d", id)
        }
        return nil, errors.Wrap(err, "query user failed")
    }
    return &user, nil
}

// 列表查询
func (r *UserRepository) FindAll(ctx context.Context) ([]*model.User, error) {
    var users []*model.User
    err := r.db.WithContext(ctx).Order("id DESC").Find(&users).Error
    if err != nil {
        return nil, errors.Wrap(err, "query users failed")
    }
    return users, nil
}

// 条件查询
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
    var user model.User
    err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.ErrNotFound
        }
        return nil, errors.Wrap(err, "query user by email failed")
    }
    return &user, nil
}

// 分页查询
func (r *UserRepository) FindPage(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
    var users []*model.User
    var total int64
    
    // 查询总数
    if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
        return nil, 0, errors.Wrap(err, "count users failed")
    }
    
    // 查询分页数据
    offset := (page - 1) * pageSize
    err := r.db.WithContext(ctx).
        Offset(offset).
        Limit(pageSize).
        Order("id DESC").
        Find(&users).Error
    
    if err != nil {
        return nil, 0, errors.Wrap(err, "query users page failed")
    }
    
    return users, total, nil
}
```

### 5. 创建方法

```go
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
    err := r.db.WithContext(ctx).Create(user).Error
    if err != nil {
        return errors.Wrap(err, "create user failed")
    }
    return nil
}
```

### 6. 更新方法

```go
// 全量更新
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
    err := r.db.WithContext(ctx).Save(user).Error
    if err != nil {
        return errors.Wrap(err, "update user failed")
    }
    return nil
}

// 部分字段更新
func (r *UserRepository) UpdateStatus(ctx context.Context, id uint, status int) error {
    err := r.db.WithContext(ctx).
        Model(&model.User{}).
        Where("id = ?", id).
        Update("status", status).Error
    if err != nil {
        return errors.Wrapf(err, "update user status failed, id: %d", id)
    }
    return nil
}
```

### 7. 删除方法

```go
// 物理删除
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
    err := r.db.WithContext(ctx).Delete(&model.User{}, id).Error
    if err != nil {
        return errors.Wrapf(err, "delete user failed, id: %d", id)
    }
    return nil
}
```

## 💡 最佳实践

### 1. 始终使用 Context

```go
// ✅ 正确
r.db.WithContext(ctx).Find(&users)

// ❌ 错误
r.db.Find(&users)
```

### 2. 错误处理

```go
// ✅ 区分 NotFound 错误
if err == gorm.ErrRecordNotFound {
    return nil, errors.ErrNotFound
}

// ✅ 包装其他错误
return nil, errors.Wrap(err, "query failed")
```

### 3. 使用查询条件

```go
// ✅ 使用链式调用
r.db.WithContext(ctx).
    Where("status = ?", 1).
    Where("created_at > ?", startTime).
    Order("id DESC").
    Limit(10).
    Find(&users)
```

### 4. 避免 N+1 查询

```go
// ✅ 预加载关联
r.db.WithContext(ctx).
    Preload("Orders").
    Preload("Profile").
    Find(&users)

// ❌ 循环查询（N+1问题）
for _, user := range users {
    r.db.Find(&user.Orders)  // 不好
}
```

### 5. 事务支持

```go
// 支持外部事务
func (r *UserRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, user *model.User) error {
    return tx.WithContext(ctx).Create(user).Error
}

// Service 层使用
func (s *UserService) CreateUserAndProfile(ctx context.Context, user *model.User) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        if err := s.userRepo.CreateWithTx(ctx, tx, user); err != nil {
            return err
        }
        // ...
        return nil
    })
}
```

## 📋 方法命名规范

### 查询方法

- `FindByID` - 根据 ID 查询单条
- `FindByXXX` - 根据条件查询单条
- `FindAll` - 查询所有
- `FindPage` - 分页查询
- `Count` - 统计数量
- `Exists` - 判断是否存在

### 修改方法

- `Create` - 创建
- `Update` - 更新
- `Delete` - 删除
- `UpdateXXX` - 更新指定字段

## 🎨 代码组织

```go
package repository

import (...)

// 1. 类型定义
type UserRepository struct {
    db *gorm.DB
}

// 2. 构造函数
func NewUserRepository(db *gorm.DB) *UserRepository {}

// 3. 查询方法
func (r *UserRepository) FindByID(...) {}
func (r *UserRepository) FindAll(...) {}
func (r *UserRepository) FindPage(...) {}

// 4. 创建方法
func (r *UserRepository) Create(...) {}

// 5. 更新方法
func (r *UserRepository) Update(...) {}

// 6. 删除方法
func (r *UserRepository) Delete(...) {}

// 7. 其他方法
func (r *UserRepository) Count(...) {}
func (r *UserRepository) Exists(...) {}
```

## 🔗 依赖注入

在 `cmd/server/wire.go` 中注册：

```go
wire.Build(
    // ...
    database.NewMySQLDB,
    repository.NewUserRepository,  // 添加这里
    // ...
)
```

## 📚 相关文档

- [Service 层说明](../service/README.md)
- [Model 层说明](../model/README.md)
- [项目结构说明](../../STRUCTURE.md)

---

**Repository 层是数据访问的唯一入口！** 💾
