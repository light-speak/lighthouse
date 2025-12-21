<p align="center">
  <img src="docs/public/logo.svg" width="80" height="80" alt="Lighthouse">
</p>

<h1 align="center">Lighthouse</h1>

<p align="center">
  <strong>用 Go 构建 GraphQL API，快。</strong>
</p>

<p align="center">
  <a href="https://github.com/light-speak/lighthouse/releases"><img src="https://img.shields.io/github/v/release/light-speak/lighthouse?style=flat-square&color=blue" alt="Release"></a>
  <a href="https://pkg.go.dev/github.com/light-speak/lighthouse"><img src="https://pkg.go.dev/badge/github.com/light-speak/lighthouse.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/light-speak/lighthouse"><img src="https://goreportcard.com/badge/github.com/light-speak/lighthouse" alt="Go Report Card"></a>
  <a href="https://github.com/light-speak/lighthouse/blob/main/LICENSE"><img src="https://img.shields.io/github/license/light-speak/lighthouse?style=flat-square" alt="License"></a>
</p>

<p align="center">
  <a href="https://light-speak.github.io/lighthouse/">在线文档</a> •
  <a href="https://light-speak.github.io/lighthouse/guide/getting-started">快速开始</a> •
  <a href="./README.md">English</a>
</p>

---

## Lighthouse 是什么？

Lighthouse 是一个**开箱即用**的 Go GraphQL 框架。定义 Schema，运行一条命令，即可获得带有 DataLoader、认证、数据库迁移的生产级 API。

```graphql
type User @loader(keys: ["id"]) {
  id: ID!
  name: String! @varchar(length: 100)
  posts: [Post!]! @goField(forceResolver: true)
}

extend type Query {
  me: User! @auth
}
```

```bash
lighthouse generate:schema  # 搞定。模型、解析器、DataLoader 全部生成。
```

## 功能特性

| 特性 | 说明 |
|------|------|
| **Schema 优先** | 定义 GraphQL Schema，自动生成 Go 代码 |
| **DataLoader** | 自动生成，N+1 问题一键解决 |
| **认证指令** | `@auth`、`@own` 开箱即用 |
| **数据库** | GORM + MySQL，连接池，主从分离 |
| **迁移** | Atlas 驱动的数据库迁移 |
| **队列** | Redis 异步任务 (asynq) |
| **消息** | NATS 发布/订阅，实时通信 |
| **存储** | S3/MinIO/COS 统一接口 |
| **监控** | Prometheus 指标 + 健康检查 |
| **MCP** | AI 辅助开发，支持 Claude Code |

## 5 分钟上手

```bash
# 安装
go install github.com/light-speak/lighthouse@latest

# 创建项目
lighthouse generate:init --module github.com/you/myapp --models user,post
cd myapp

# 配置 .env，然后运行
go run . app:start
```

打开 http://localhost:8080 → GraphQL Playground 就绪。

## 项目结构

```
myapp/
├── schema/          # GraphQL 定义
├── models/          # 生成的 Go 结构体
├── resolver/        # 你的业务逻辑
├── graph/           # gqlgen 生成（别动）
├── commands/        # CLI 命令
├── server/          # HTTP 服务配置
└── migrations/      # 数据库迁移
```

## 技术栈

- **[gqlgen](https://gqlgen.com/)** - GraphQL 引擎
- **[GORM](https://gorm.io/)** - ORM
- **[Atlas](https://atlasgo.io/)** - 数据库迁移
- **[asynq](https://github.com/hibiken/asynq)** - 任务队列
- **[NATS](https://nats.io/)** - 消息系统
- **[zerolog](https://github.com/rs/zerolog)** - 日志

## MCP 集成

Lighthouse 内置 [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) 支持，让 Claude Code 等 AI 助手能够理解并协助开发 Lighthouse 项目。

### 配置

**方式一：命令行配置（推荐）**

```bash
# 先安装 lighthouse
go install github.com/light-speak/lighthouse@latest

# 添加 MCP 服务器到 Claude Code
claude mcp add lighthouse -- lighthouse mcp
```

**方式二：手动配置**

在项目根目录创建 `.mcp.json`（团队共享）：

```json
{
  "mcpServers": {
    "lighthouse": {
      "type": "stdio",
      "command": "lighthouse",
      "args": ["mcp"]
    }
  }
}
```

或添加到 `~/.claude.json`（个人跨项目使用）：

```json
{
  "mcpServers": {
    "lighthouse": {
      "command": "lighthouse",
      "args": ["mcp"]
    }
  }
}
```

**验证安装**

```bash
# 检查 MCP 服务器是否已注册
claude mcp list
```

### AI 能做什么

| 能力 | 说明 |
|------|------|
| **生成代码** | 创建 schema、resolver、命令、任务 |
| **查询文档** | 获取指令用法（@loader, @auth 等） |
| **配置帮助** | 数据库、Redis、队列配置说明 |
| **执行命令** | 运行任意 lighthouse CLI 命令 |
| **读取示例** | 访问 schema、resolver、dataloader 示例 |

### 可用工具

```
generate_schema      - 生成 GraphQL schema 和 models
generate_dataloader  - 生成 DataLoader 代码
generate_command     - 创建 CLI 命令
generate_task        - 创建异步队列任务
init_project         - 初始化新项目
get_directive_info   - 获取指令文档
get_config_info      - 获取配置文档
list_generators      - 列出所有生成器
search_docs          - 搜索框架文档
run_command          - 运行任意 lighthouse 命令
```

## 文档

📚 **[完整文档](https://light-speak.github.io/lighthouse/)**

- [快速开始](https://light-speak.github.io/lighthouse/guide/getting-started)
- [Schema 基础](https://light-speak.github.io/lighthouse/schema/basics)
- [DataLoader](https://light-speak.github.io/lighthouse/schema/dataloader)
- [认证中间件](https://light-speak.github.io/lighthouse/features/auth)
- [数据库](https://light-speak.github.io/lighthouse/features/database)

## 贡献

欢迎 PR。重大改动请先开 Issue 讨论。

## 许可证

MIT

---

<p align="center">
  <sub>由 <a href="https://github.com/light-speak">Light Speak</a> 用 ☕ 驱动</sub>
</p>
