# Go API Template

> 基于 Gin 的企业级 Go API 项目模板，开箱即用。

## ✨ 特性

- 🏗️ **清晰的分层架构** - Controller → Service → Repository
- 💉 **依赖注入** - Google Wire 自动生成代码
- 🔌 **框架隔离** - 完全隔离 Gin 依赖，易于迁移
- 📝 **企业级日志** - Zap + Lumberjack 高性能日志
- ⚠️ **完整错误处理** - cockroachdb/errors 堆栈跟踪
- 💾 **缓存支持** - Redis/Memory/Chain 多种驱动
- 🗄️ **数据库 ORM** - GORM 支持 MySQL/PostgreSQL
- 📦 **配置管理** - YAML 配置文件
- 🎯 **Demo 示例** - 完整的 CRUD 实现
- 🌐 **CORS 支持** - 可配置的跨域处理（已集成）

## 🚀 快速开始

### 环境要求

- Go 1.21+
- MySQL 5.7+ / PostgreSQL 12+ (可选)
- Redis 6.0+ (可选，默认使用内存缓存)

### 方式一：使用脚本创建（推荐⭐）

一键创建新项目，自动处理所有配置：

```bash
cd go-api-template
./create-project.sh
```

**交互式输入：**

```
📝 请输入项目名称: my-api

📦 请输入 Go 模块路径: github.com/username/my-api
# 或简单名称: my-api
# 或企业项目: company.com/my-api

确认创建？(y/N): y
```

**脚本自动完成：**
- ✅ 在**同级目录**创建新项目
- ✅ 复制完整的代码和配置
- ✅ 替换模块路径
- ✅ 安装所有依赖
- ✅ 自动处理域名格式（添加 replace 指令）
- ✅ 生成项目文档

**创建后的目录：**

```bash
Projects/golang/
├── go-api-template/    # 模板（保持不变）
└── my-api/             # 新项目（同级目录）
```

### 方式二：手动创建（高级用户）

```bash
# 1. 复制模板
cp -r go-api-template my-project
cd my-project
rm create-project.sh  # 删除创建脚本

# 2. 修改 go.mod
vim go.mod
# 将 module go-api-template 改为 module your-project

# 3. 全局替换导入路径
find . -type f -name "*.go" -exec sed -i '' 's|go-api-template|your-project|g' {} +

# 4. 安装依赖
go mod tidy
go install github.com/google/wire/cmd/wire@latest
```

### 启动项目

```bash
# 进入项目目录（如果使用脚本创建）
cd ../my-api

# 1. 生成 Wire 依赖注入代码
make wire

# 2. 运行项目
make run

# 启动成功后会显示：
# 🌐 服务地址: http://localhost:8080
# 📚 API 端点: ...
```

### 配置数据库（可选）

Demo 示例需要 MySQL 数据库。如果要测试完整的 CRUD 功能，需要先创建数据库和表：

```sql
-- 1. 创建数据库
CREATE DATABASE go_api_template DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 2. 使用数据库
USE go_api_template;

-- 3. 创建 demos 表
CREATE TABLE `demos` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `title` varchar(200) NOT NULL COMMENT '标题',
  `content` text COMMENT '内容',
  `status` int NOT NULL DEFAULT '1' COMMENT '状态：1-启用 0-禁用',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Demo示例表';

-- 4. 插入测试数据（可选）
INSERT INTO `demos` (`title`, `content`, `status`, `created_at`, `updated_at`) VALUES
('第一个Demo', '这是第一个演示内容', 1, NOW(), NOW()),
('第二个Demo', '这是第二个演示内容', 1, NOW(), NOW()),
('第三个Demo', '这是第三个演示内容', 0, NOW(), NOW());
```

然后修改配置文件 `config/config.yaml`：

```yaml
database:
  driver: mysql
  host: localhost
  port: 3306
  username: root          # 你的数据库用户名
  password: password      # 你的数据库密码
  database: go_api_template
```

### 测试 API

```bash
# 健康检查（无需数据库）
curl http://localhost:8080/health

# 获取所有 Demo
curl http://localhost:8080/api/v1/demos

# 创建 Demo
curl -X POST http://localhost:8080/api/v1/demos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "测试标题",
    "content": "测试内容",
    "status": 1
  }'

# 获取单个 Demo（假设 ID 为 1）
curl http://localhost:8080/api/v1/demos/1

# 更新 Demo
curl -X PUT http://localhost:8080/api/v1/demos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "更新后的标题",
    "content": "更新后的内容",
    "status": 1
  }'

# 删除 Demo
curl -X DELETE http://localhost:8080/api/v1/demos/1
```

**响应示例：**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "title": "测试标题",
    "content": "测试内容",
    "status": 1,
    "created_at": "2024-01-01T12:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

### 脚本执行流程

1. **检查环境** - 验证 Go 是否安装
2. **获取项目信息** - 交互式输入项目名称和模块路径
3. **创建目录** - 在模板同级目录创建新项目
4. **复制文件** - 复制所有模板代码和配置
5. **初始化模块** - 使用临时名称初始化（避免网络问题）
6. **安装依赖** - 安装所有第三方依赖
7. **替换模块名** - 将临时名称替换为目标模块路径
8. **配置本地开发** - 域名格式自动添加 replace 指令

### 模块路径选择建议

```bash
# ✅ 推荐：简单名称（本地/内部项目）
my-api
user-service

# ✅ 推荐：GitHub（开源项目）
github.com/username/my-api

# ✅ 支持：企业域名（自动配置 replace）
company.com/team/my-api
gitlab.company.com/backend/my-api

# ⚠️  注意：域名格式会自动添加 replace 指令，避免网络下载
```

### 使用示例

```bash
$ cd go-api-template
$ ./create-project.sh

========================================
   Go API Template - 项目创建工具
========================================

ℹ️  检查依赖...
✅ 依赖检查通过

ℹ️  项目名称：用作项目目录名
📝 请输入项目名称: user-service

ℹ️  Go 模块路径：go.mod 中的 module 声明
📦 请输入 Go 模块路径: github.com/myname/user-service

ℹ️  项目信息确认：
   项目名称: user-service
   模块路径: github.com/myname/user-service
   创建位置: /Users/you/Projects/golang/user-service

确认创建？(y/N): y

ℹ️  创建项目目录...
ℹ️  创建目录结构...
ℹ️  复制模板文件...
ℹ️  初始化 Go 模块...
ℹ️  配置导入路径...
ℹ️  创建应用入口...
ℹ️  安装依赖...
ℹ️  更新模块路径: github.com/myname/user-service
✅ 模块路径已更新
ℹ️  检测到域名格式，添加本地模块配置...
✅ 已添加 replace 指令
ℹ️  创建项目文档...

✅ 项目创建成功！

========================================
📁 项目位置: /Users/you/Projects/golang/user-service
========================================

ℹ️  下一步操作：
   1. cd /Users/you/Projects/golang/user-service
   2. make wire
   3. make run
   4. curl http://localhost:8080/health

✅ 开始构建你的 API 项目！🚀
```

---

## 📚 API 示例

模板包含完整的 Demo CRUD 示例：

```bash
GET    /api/v1/demos       # 获取所有 Demo
GET    /api/v1/demos/:id   # 获取单个 Demo
POST   /api/v1/demos       # 创建 Demo
PUT    /api/v1/demos/:id   # 更新 Demo
DELETE /api/v1/demos/:id   # 删除 Demo
```

## 🛠️ 开发

### 常用命令

```bash
make wire          # 生成 Wire 依赖注入代码
make build         # 编译项目
make run           # 运行项目
make test          # 运行测试
make clean         # 清理编译文件
make help          # 查看所有命令
```

### 添加新功能

1. 在 `internal/model/` 定义数据模型
2. 在 `internal/repository/` 实现数据访问
3. 在 `internal/service/` 实现业务逻辑
4. 在 `internal/controller/` 实现 HTTP 接口
5. 在 `cmd/server/wire.go` 注册依赖
6. 运行 `make wire` 生成代码

详细说明请参考 `internal/controller/README.md`

## 📖 文档

### 项目级文档
- [项目结构说明](./STRUCTURE.md) - 详细的目录结构和架构设计

### 分层架构开发指南
- [Controller 层](./internal/controller/README.md) - HTTP 控制器开发
- [Service 层](./internal/service/README.md) - 业务逻辑开发
- [Repository 层](./internal/repository/README.md) - 数据访问开发
- [Model 层](./internal/model/README.md) - 数据模型定义
- [Middleware 层](./internal/middleware/README.md) - 中间件开发
- [Constants 包](./internal/constants/README.md) - 常量管理

## ⚙️ 配置

主要配置项（`config/config.yaml`）：

```yaml
server:
  port: 8080              # 服务端口
  mode: debug             # debug, release, test

database:
  driver: mysql           # mysql, postgres
  host: localhost
  port: 3306
  username: root
  password: password
  database: go_api_template  # 数据库名（需先创建）

redis:
  host: localhost
  port: 6379
  password: ""            # Redis 密码（如有）

cache:
  driver: memory          # redis, memory, chain
  ttl: 300                # 默认过期时间（秒）

logger:
  level: info             # debug, info, warn, error
  filename: logs/app.log
  console: true           # 是否输出到控制台

cors:
  enabled: true           # 是否启用 CORS
  allow_origins:          # 允许的来源
    - "*"                 # "*" 允许所有来源
  allow_methods:          # 允许的 HTTP 方法
    - "GET"
    - "POST"
    - "PUT"
    - "DELETE"
  allow_headers:          # 允许的请求头
    - "Content-Type"
    - "Authorization"
```

**数据库设置：**

模板默认使用**内存缓存**，可以在不配置数据库的情况下运行（但 Demo CRUD API 需要数据库）。

如果要使用完整的 Demo 示例：
1. 创建数据库（见上方"配置数据库"章节的 DDL）
2. 修改 `config/config.yaml` 中的数据库配置
3. 重启服务

## 🏗️ 架构设计

### 分层架构

```
┌─────────────────────────────────────┐
│         HTTP Request                │
└─────────────┬───────────────────────┘
              │
      ┌───────▼────────┐
      │   Controller   │  ← HTTP 处理、参数验证
      └───────┬────────┘
              │
      ┌───────▼────────┐
      │    Service     │  ← 业务逻辑、事务控制
      └───────┬────────┘
              │
      ┌───────▼────────┐
      │   Repository   │  ← 数据访问、数据库操作
      └───────┬────────┘
              │
      ┌───────▼────────┐
      │    Database    │
      └────────────────┘
```

### 依赖注入

使用 Google Wire 实现编译时依赖注入：
- 自动生成依赖关系
- 编译时检查
- 无运行时反射开销

### 框架隔离

通过 `pkg/web` 隔离 Gin 框架：
- 业务代码不直接依赖 Gin
- 易于测试和框架迁移
- 降低耦合度

## 📦 核心依赖

- [gin-gonic/gin](https://github.com/gin-gonic/gin) - Web 框架
- [google/wire](https://github.com/google/wire) - 依赖注入
- [uber-go/zap](https://github.com/uber-go/zap) - 日志库
- [gorm.io/gorm](https://gorm.io/) - ORM 框架
- [cockroachdb/errors](https://github.com/cockroachdb/errors) - 错误处理

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 License

MIT License

---

**开始使用这个模板构建你的 Go API 项目！** 🚀
