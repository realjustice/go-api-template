# Database 包说明

## 📖 包的作用

提供数据库连接和基础数据访问功能。

## 📁 文件说明

- `mysql.go` - MySQL 数据库连接
- `base_repository.go` - 基础 Repository，提供通用 CRUD 操作

## 🎯 BaseRepository - 通用数据访问

`BaseRepository` 封装了最常用的数据库操作，减少重复代码。

### 设计理念

**部分封装，保留灵活性：**
- ✅ 封装通用的 CRUD 操作（80% 场景）
- ✅ 复杂查询直接使用 GORM（20% 场景）
- ✅ 通过嵌入使用，不是强制的
- ✅ 可以随时访问底层 GORM

### 使用方式

#### 1. 嵌入 BaseRepository

```go
package repository

import (
    "go-api-template/pkg/database"
    "gorm.io/gorm"
)

type UserRepository struct {
    *database.BaseRepository  // 嵌入基类
    db *gorm.DB              // 保留 db 用于复杂查询
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{
        BaseRepository: database.NewBaseRepository(db),
        db:             db,
    }
}
```

#### 2. 使用基类方法（简单操作）

```go
// 根据 ID 查询
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
    var user model.User
    err := r.BaseRepository.FindByID(ctx, id, &user)  // 使用基类方法
    return &user, err
}

// 创建
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
    return r.BaseRepository.Create(ctx, user)  // 使用基类方法
}

// 更新
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
    return r.BaseRepository.Update(ctx, user)  // 使用基类方法
}

// 删除
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
    return r.BaseRepository.Delete(ctx, &model.User{}, id)  // 使用基类方法
}
```

#### 3. 直接使用 GORM（复杂查询）

```go
// 复杂查询，直接使用 GORM
func (r *UserRepository) SearchUsers(ctx context.Context, keyword string, status int) ([]*model.User, error) {
    var users []*model.User
    
    // 直接使用 r.db，保留 GORM 的全部灵活性
    err := r.db.WithContext(ctx).
        Where("name LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%").
        Where("status = ?", status).
        Preload("Profile").       // 预加载关联
        Preload("Orders").        // 预加载关联
        Order("created_at DESC").
        Limit(100).
        Find(&users).Error
    
    return users, err
}
```

#### 4. 混合使用

```go
// 简单查询用基类
func (r *UserRepository) FindAll(ctx context.Context) ([]*model.User, error) {
    var users []*model.User
    err := r.BaseRepository.FindAll(ctx, &users, "1 = 1")
    return users, err
}

// 复杂查询用 GORM
func (r *UserRepository) FindActiveUsersWithOrders(ctx context.Context) ([]*model.User, error) {
    var users []*model.User
    err := r.db.WithContext(ctx).
        Where("status = ?", 1).
        Preload("Orders", "status = ?", "paid").
        Find(&users).Error
    return users, err
}
```

## 📋 BaseRepository 提供的方法

### 查询方法

| 方法 | 说明 | 使用场景 |
|------|------|----------|
| `FindByID` | 根据 ID 查询 | 查询单条记录 |
| `FindOne` | 根据条件查询单条 | 查询单条记录 |
| `FindAll` | 查询所有 | 列表查询 |
| `FindPage` | 分页查询 | 分页列表 |
| `Count` | 统计数量 | 统计 |
| `Exists` | 判断是否存在 | 验证 |

### 创建方法

| 方法 | 说明 |
|------|------|
| `Create` | 创建单条 |
| `CreateInBatches` | 批量创建 |

### 更新方法

| 方法 | 说明 |
|------|------|
| `Update` | 更新全部字段 |
| `UpdateFields` | 更新指定字段 |
| `UpdateColumn` | 更新单个字段 |

### 删除方法

| 方法 | 说明 |
|------|------|
| `Delete` | 根据 ID 删除 |
| `DeleteWhere` | 根据条件删除 |

### 事务和 SQL

| 方法 | 说明 |
|------|------|
| `Transaction` | 执行事务 |
| `Exec` | 执行原生 SQL |
| `Raw` | 原生查询 |
| `DB` | 获取 GORM 实例 |

## 💡 使用示例

### 示例 1：简单 CRUD（使用 BaseRepository）

```go
type ProductRepository struct {
    *database.BaseRepository
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
    return &ProductRepository{
        BaseRepository: database.NewBaseRepository(db),
        db:             db,
    }
}

// ✅ 使用基类方法 - 简洁
func (r *ProductRepository) FindByID(ctx context.Context, id uint) (*model.Product, error) {
    var product model.Product
    err := r.BaseRepository.FindByID(ctx, id, &product)
    return &product, err
}

func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
    return r.BaseRepository.Create(ctx, product)
}
```

### 示例 2：复杂查询（直接使用 GORM）

```go
// ✅ 直接使用 GORM - 灵活
func (r *ProductRepository) FindPopularProducts(ctx context.Context, limit int) ([]*model.Product, error) {
    var products []*model.Product
    
    err := r.db.WithContext(ctx).
        Where("status = ?", 1).
        Where("sales > ?", 100).
        Preload("Category").
        Preload("Reviews", func(db *gorm.DB) *gorm.DB {
            return db.Where("rating >= ?", 4).Order("created_at DESC")
        }).
        Order("sales DESC").
        Limit(limit).
        Find(&products).Error
    
    return products, err
}
```

### 示例 3：分页查询

```go
func (r *ProductRepository) FindPage(ctx context.Context, page, pageSize int) ([]*model.Product, int64, error) {
    var products []*model.Product
    total, err := r.BaseRepository.FindPage(ctx, &products, page, pageSize, "status = ?", 1)
    return products, total, err
}
```

### 示例 4：事务操作

```go
func (r *ProductRepository) CreateWithInventory(ctx context.Context, product *model.Product, inventory *model.Inventory) error {
    return r.BaseRepository.Transaction(ctx, func(tx *gorm.DB) error {
        // 创建产品
        if err := tx.Create(product).Error; err != nil {
            return err
        }
        
        // 创建库存
        inventory.ProductID = product.ID
        if err := tx.Create(inventory).Error; err != nil {
            return err
        }
        
        return nil
    })
}
```

### 示例 5：获取 GORM 实例

```go
func (r *ProductRepository) ComplexQuery(ctx context.Context) ([]*model.Product, error) {
    var products []*model.Product
    
    // 需要完全的 GORM 控制时
    db := r.BaseRepository.DB(ctx)  // 获取 GORM 实例
    
    err := db.
        Joins("LEFT JOIN categories ON products.category_id = categories.id").
        Where("categories.status = ?", 1).
        Group("products.id").
        Having("COUNT(*) > ?", 5).
        Find(&products).Error
    
    return products, err
}
```

## 🎯 何时使用 BaseRepository？何时使用 GORM？

### 使用 BaseRepository（推荐场景）

```go
// ✅ 简单的 CRUD
FindByID(ctx, id)
Create(ctx, model)
Update(ctx, model)
Delete(ctx, id)

// ✅ 简单的条件查询
FindOne(ctx, dest, "email = ?", email)
FindAll(ctx, dest, "status = ?", 1)

// ✅ 分页查询
FindPage(ctx, dest, page, pageSize, query, args...)

// ✅ 统计和判断
Count(ctx, model, query, args...)
Exists(ctx, model, query, args...)
```

### 直接使用 GORM（推荐场景）

```go
// ✅ 复杂查询
db.Where(...).Where(...).Preload(...).Joins(...).Find(...)

// ✅ 子查询
db.Where("id IN (?)", db.Table("orders").Select("user_id"))

// ✅ 原生 SQL
db.Raw("SELECT ... FROM ... WHERE ...").Scan(...)

// ✅ 批量操作
db.Where("status = ?", 0).Delete(&User{})

// ✅ 关联查询
db.Preload("Orders").Preload("Profile").Find(&users)
```

## 💡 最佳实践

### 1. 优先使用 BaseRepository

```go
// ✅ 简单操作，用基类
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
    var user model.User
    err := r.BaseRepository.FindByID(ctx, id, &user)
    return &user, err
}
```

### 2. 复杂查询用 GORM

```go
// ✅ 复杂查询，用 GORM
func (r *UserRepository) SearchUsers(ctx context.Context, params SearchParams) ([]*model.User, error) {
    query := r.db.WithContext(ctx)
    
    if params.Keyword != "" {
        query = query.Where("name LIKE ?", "%"+params.Keyword+"%")
    }
    
    if len(params.Tags) > 0 {
        query = query.Where("tags IN ?", params.Tags)
    }
    
    var users []*model.User
    err := query.Preload("Profile").Find(&users).Error
    return users, err
}
```

### 3. 不要过度抽象

```go
// ❌ 不推荐：为每个简单方法都定义接口
type UserRepository interface {
    FindByID(ctx context.Context, id uint) (*model.User, error)
    Create(ctx context.Context, user *model.User) error
    // ...
}

// ✅ 推荐：直接使用具体类型
type UserRepository struct {
    *database.BaseRepository
    db *gorm.DB
}
```

## 🔄 迁移到其他 ORM

如果将来真的需要换 ORM，只需要：

### 方案 1：重写 BaseRepository

```go
// 新的 BaseRepository 实现（使用其他 ORM）
type BaseRepository struct {
    db *sqlx.DB  // 改用 sqlx
}

// 重新实现方法
func (r *BaseRepository) FindByID(ctx context.Context, id interface{}, dest interface{}) error {
    // 使用 sqlx 实现
}
```

### 方案 2：重写各个 Repository

```go
// 只需要改 Repository 层的实现
type UserRepository struct {
    db *sqlx.DB  // 改用其他 ORM
}

// Service 和 Controller 完全不受影响
```

## 📚 相关文档

- [Repository 层说明](../../internal/repository/README.md)
- [项目结构说明](../../STRUCTURE.md)

---

**BaseRepository 让 80% 的数据访问代码更简洁！** 🚀
