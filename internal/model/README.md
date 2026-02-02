# Model 层说明

## 📖 包的作用

Model 层定义数据模型，对应数据库表结构，用于数据的存储和传输。

## 🎯 职责范围

### ✅ Model 层应该做什么

- 定义数据结构
- 配置 GORM 标签（表映射）
- 配置 JSON 标签（序列化）
- 定义表关联关系
- 定义表名映射
- 定义数据约束

### ❌ Model 层不应该做什么

- 不包含业务逻辑
- 不包含数据库操作
- 不包含验证逻辑（在 Service 或 Controller）
- 保持纯粹的数据定义

## 📝 示例代码

参考 `demo.go`，这是一个标准的 Model 定义。

### 基本结构

```go
package model

import "time"

// Demo 演示模型
type Demo struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Title     string    `json:"title" gorm:"type:varchar(200);not null"`
    Content   string    `json:"content" gorm:"type:text"`
    Status    int       `json:"status" gorm:"default:1;comment:状态 1-启用 0-禁用"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Demo) TableName() string {
    return "demos"
}
```

## 🏗️ 开发新的 Model

### 1. 文件命名

```
{模块名}.go

示例：
user.go
order.go
product.go
```

### 2. 完整示例

```go
package model

import "time"

// User 用户模型
type User struct {
    // 主键
    ID uint `json:"id" gorm:"primaryKey;comment:用户ID"`
    
    // 基本字段
    Username string `json:"username" gorm:"type:varchar(50);not null;uniqueIndex;comment:用户名"`
    Email    string `json:"email" gorm:"type:varchar(100);not null;uniqueIndex;comment:邮箱"`
    Phone    string `json:"phone" gorm:"type:varchar(20);index;comment:手机号"`
    Password string `json:"-" gorm:"type:varchar(255);not null;comment:密码"`  // json:"-" 不序列化
    
    // 状态字段
    Status int `json:"status" gorm:"type:tinyint;default:1;index;comment:状态 1-正常 0-禁用"`
    
    // 时间字段
    CreatedAt time.Time  `json:"created_at" gorm:"comment:创建时间"`
    UpdatedAt time.Time  `json:"updated_at" gorm:"comment:更新时间"`
    DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index;comment:删除时间"` // 软删除
}

// TableName 指定表名
func (User) TableName() string {
    return "users"
}
```

## 🏷️ GORM 标签说明

### 字段类型

```go
gorm:"type:varchar(100)"      // 字符串类型
gorm:"type:text"              // 文本类型
gorm:"type:int"               // 整数类型
gorm:"type:bigint"            // 大整数
gorm:"type:decimal(10,2)"     // 小数
gorm:"type:datetime"          // 日期时间
gorm:"type:json"              // JSON 类型
```

### 约束

```go
gorm:"primaryKey"             // 主键
gorm:"not null"               // 非空
gorm:"unique"                 // 唯一约束
gorm:"default:1"              // 默认值
gorm:"autoIncrement"          // 自增
```

### 索引

```go
gorm:"index"                  // 普通索引
gorm:"uniqueIndex"            // 唯一索引
gorm:"index:idx_name"         // 命名索引
gorm:"index:,composite:true"  // 复合索引
```

### 其他

```go
gorm:"comment:字段说明"       // 字段注释
gorm:"column:col_name"        // 自定义列名
gorm:"-"                      // 忽略该字段
gorm:"embedded"               // 嵌入结构体
```

## 🎨 JSON 标签说明

```go
json:"field_name"             // JSON 字段名
json:"-"                      // 不序列化（如：密码）
json:"field,omitempty"        // 空值时忽略
```

## 💡 最佳实践

### 1. 基础字段模板

```go
type BaseModel struct {
    ID        uint       `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`  // 软删除
}

// 使用基础模板
type User struct {
    BaseModel
    Username string `json:"username" gorm:"type:varchar(50);not null;uniqueIndex"`
    Email    string `json:"email" gorm:"type:varchar(100);not null;uniqueIndex"`
}
```

### 2. 敏感字段处理

```go
type User struct {
    Password string `json:"-" gorm:"type:varchar(255);not null"`  // 不返回给前端
    Salt     string `json:"-" gorm:"type:varchar(64)"`            // 不返回给前端
}
```

### 3. 关联关系

```go
// 一对多
type User struct {
    ID     uint     `json:"id" gorm:"primaryKey"`
    Orders []*Order `json:"orders,omitempty" gorm:"foreignKey:UserID"`
}

type Order struct {
    ID     uint  `json:"id" gorm:"primaryKey"`
    UserID uint  `json:"user_id" gorm:"index"`
    User   *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// 多对多
type User struct {
    ID    uint    `json:"id" gorm:"primaryKey"`
    Roles []*Role `json:"roles,omitempty" gorm:"many2many:user_roles"`
}
```

### 4. 自定义表名

```go
// 方法 1：实现 TableName 方法（推荐）
func (User) TableName() string {
    return "users"  // 自定义表名
}

// 方法 2：使用标签
type User struct {
    // ...
} // `gorm:"table:users"`  // 不推荐，用方法更灵活
```

### 5. 默认值和自动时间

```go
type User struct {
    // 数据库默认值
    Status int `gorm:"default:1"`
    
    // GORM 自动管理时间
    CreatedAt time.Time  // 创建时自动设置
    UpdatedAt time.Time  // 更新时自动更新
    DeletedAt *time.Time `gorm:"index"`  // 软删除时自动设置
}
```

### 6. 复合索引

```go
type User struct {
    Email  string `gorm:"index:idx_email_status,priority:1"`
    Status int    `gorm:"index:idx_email_status,priority:2"`
    // 会创建复合索引：idx_email_status (email, status)
}
```

## 📋 命名规范

### 结构体命名

```go
type User struct {}      // 单数，首字母大写
type Order struct {}
type Product struct {}
```

### 字段命名

```go
type User struct {
    ID       uint      // 缩写大写
    UserID   uint      // 复合词，每个单词首字母大写
    Username string    // 复合词，驼峰命名
    Email    string    // 单词首字母大写
}
```

### 表名

```go
func (User) TableName() string {
    return "users"      // 复数，全小写，下划线分隔
}

func (OrderItem) TableName() string {
    return "order_items"  // 复合词用下划线
}
```

## 🎨 代码组织

```go
package model

import (
    "time"
)

// 1. 主模型定义
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username" gorm:"type:varchar(50);not null"`
    CreatedAt time.Time `json:"created_at"`
}

// 2. 表名方法
func (User) TableName() string {
    return "users"
}

// 3. 辅助方法（如需要）
func (u *User) IsActive() bool {
    return u.Status == 1
}
```

## 🗄️ 对应的 DDL

为 Model 创建对应的数据库表：

```sql
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password` varchar(255) NOT NULL,
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '1-正常 0-禁用',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_email` (`email`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## 📚 相关文档

- [Repository 层说明](../repository/README.md)
- [Service 层说明](../service/README.md)
- [项目结构说明](../../STRUCTURE.md)
- [GORM 官方文档](https://gorm.io/zh_CN/docs/)

---

**Model 层是数据的定义和映射！** 🗂️
