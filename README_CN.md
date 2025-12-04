<p align="center">
  <h1 align="center">Lighthouse</h1>
  <p align="center">功能完备的 Go GraphQL 框架</p>
</p>

<p align="center">
  <a href="https://github.com/light-speak/lighthouse/releases"><img src="https://img.shields.io/github/v/release/light-speak/lighthouse?style=flat-square" alt="Release"></a>
  <a href="https://pkg.go.dev/github.com/light-speak/lighthouse"><img src="https://pkg.go.dev/badge/github.com/light-speak/lighthouse.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/light-speak/lighthouse"><img src="https://goreportcard.com/badge/github.com/light-speak/lighthouse" alt="Go Report Card"></a>
  <a href="https://github.com/light-speak/lighthouse/blob/main/LICENSE"><img src="https://img.shields.io/github/license/light-speak/lighthouse?style=flat-square" alt="License"></a>
  <a href="https://github.com/light-speak/lighthouse/stargazers"><img src="https://img.shields.io/github/stars/light-speak/lighthouse?style=flat-square" alt="Stars"></a>
</p>

<p align="center">
  <a href="./README.md">English</a>
</p>

---

## 概述

**Lighthouse** 是一个功能完备的 Go GraphQL 框架，灵感来源于 [Laravel Lighthouse](https://lighthouse-php.com/)。它提供了一种优雅且高效的方式来构建 GraphQL API，内置数据库管理、缓存、消息队列等功能。

## 特性

- **GraphQL 优先** - 基于 [gqlgen](https://gqlgen.com/) 构建，Go 语言最流行的 GraphQL 库
- **数据库管理** - 基于 GORM 的 MySQL 支持，包含连接池、主从复制
- **DataLoader** - 内置 DataLoader 模式，高效解决 N+1 查询问题
- **Redis 集成** - 连接池、缓存工具和发布/订阅支持
- **异步队列** - 基于 [asynq](https://github.com/hibiken/asynq) 的后台任务处理
- **消息系统** - 基于 NATS 的实时消息和事件广播
- **身份认证** - 基于 JWT 的认证，支持 GraphQL 指令
- **代码生成** - CLI 工具快速生成模型、解析器和 DataLoader
- **优雅关闭** - 应用终止时正确清理资源

## 环境要求

- Go 1.24+
- MySQL 5.7+ 或 8.0+
- Redis 6.0+ (可选)
- NATS 2.0+ (可选)

## 安装

```bash
go get github.com/light-speak/lighthouse
```

### 安装 CLI 工具

```bash
go install github.com/light-speak/lighthouse@latest
```

## 快速开始

### 1. 初始化新项目

```bash
lighthouse generate:init --module=github.com/yourname/myproject --models=user,post
cd myproject
```

### 3. 配置环境变量

编辑 `.env` 文件，填写数据库配置：

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=myproject
```

### 4. 定义 GraphQL Schema

编辑 `graph/schema.graphqls`：

```graphql
type User {
  id: ID!
  name: String!
  email: String!
  createdAt: DateTime!
}

type Query {
  users: [User!]!
  user(id: ID!): User
}

type Mutation {
  createUser(name: String!, email: String!): User!
}
```

### 5. 生成代码

```bash
lighthouse generate:schema
```

### 6. 运行服务

```bash
go run . app:start
```

GraphQL 服务现在运行在 `http://localhost:8080/graphql`

## 配置说明

### 数据库连接池

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `DB_MAX_IDLE_CONNS` | 10 | 最大空闲连接数 |
| `DB_MAX_OPEN_CONNS` | 100 | 最大打开连接数 |
| `DB_CONN_MAX_LIFETIME` | 30 | 连接最大生命周期（分钟） |
| `DB_CONN_MAX_IDLE_TIME` | 3 | 空闲连接最大存活时间（分钟） |
| `DB_PREPARE_STMT` | false | 启用 prepared statement 缓存 |

### Redis 连接池

| 环境变量 | 默认值 | 说明 |
|----------|--------|------|
| `REDIS_ENABLE` | false | 是否启用 Redis |
| `REDIS_HOST` | localhost | Redis 主机 |
| `REDIS_PORT` | 6379 | Redis 端口 |
| `REDIS_PASSWORD` | | Redis 密码 |
| `REDIS_DB` | 0 | Redis 数据库 |
| `REDIS_POOL_SIZE` | 10 | 连接池大小 |
| `REDIS_MIN_IDLE_CONNS` | 5 | 最小空闲连接数 |

完整配置选项请参考 [.env.example](./lightcmd/initization/tpl/env.tpl)。

## 项目结构

```
lighthouse/
├── databases/       # 数据库连接管理 (GORM + MySQL)
├── redis/           # Redis 客户端管理
├── queue/           # 异步任务队列 (asynq)
├── messaging/       # 消息系统 (NATS)
├── routers/         # GraphQL 路由和中间件
│   ├── auth/        # 认证指令
│   ├── dataloader/  # DataLoader 批量查询
│   └── health/      # 健康检查端点
├── lightcmd/        # CLI 和代码生成
│   ├── generate/    # 代码生成器
│   └── initization/ # 项目初始化模板
├── logs/            # 日志模块
├── storages/        # 存储适配器 (S3/COS)
├── templates/       # 模板引擎
├── utils/           # 工具函数
└── lighterr/        # 错误处理
```

## CLI 命令

```bash
# 显示所有可用命令
lighthouse help

# 初始化新项目
lighthouse generate:init --module=<module-name> --models=<model1,model2>

# 生成 schema（模型、解析器、DataLoader）
lighthouse generate:schema

# 生成新命令
lighthouse generate:command --name=<command-name>

# 生成新任务
lighthouse generate:task --name=<task-name>

# 初始化队列服务
lighthouse queue:init
```

## DataLoader 模式

Lighthouse 自动生成 DataLoader 来防止 N+1 查询问题：

```go
// 自动生成的 DataLoader 用法
user, err := GetUserIdLoader(ctx).Load(ctx, userID)

// 批量加载
users, err := GetUserIdLoader(ctx).LoadAll(ctx, userIDs)
```

## 身份认证

使用 GraphQL 指令进行身份认证：

```graphql
directive @auth on FIELD_DEFINITION

type Query {
  me: User! @auth
  publicData: String!
}
```

## 健康检查

内置 Kubernetes 健康检查端点：

- **存活检查**: `GET /health`
- **就绪检查**: `GET /ready`

## 版本说明

本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

| 版本 | 状态 | Go 版本 |
|------|------|---------|
| v1.1.x | 当前版本 | 1.24+ |
| v1.0.x | 维护中 | 1.21+ |

## 贡献指南

欢迎贡献代码！请随时提交 Pull Request。

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 致谢

- [gqlgen](https://gqlgen.com/) - Go GraphQL 服务器库
- [GORM](https://gorm.io/) - Go ORM 库
- [Laravel Lighthouse](https://lighthouse-php.com/) - 本项目的灵感来源
- [asynq](https://github.com/hibiken/asynq) - 异步任务处理
- [NATS](https://nats.io/) - 消息系统

## 支持

- [文档](https://github.com/light-speak/lighthouse/wiki)
- [问题反馈](https://github.com/light-speak/lighthouse/issues)
- [讨论区](https://github.com/light-speak/lighthouse/discussions)

---

<p align="center">Made with ❤️ by <a href="https://github.com/light-speak">Light Speak</a></p>
