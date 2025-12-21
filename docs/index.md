---
layout: home

hero:
  name: Lighthouse
  text: Go GraphQL Framework
  tagline: 快速构建 GraphQL API 的现代化 Go 框架
  actions:
    - theme: brand
      text: 快速开始
      link: /guide/getting-started
    - theme: alt
      text: GitHub
      link: https://github.com/light-speak/lighthouse

features:
  - icon: 🚀
    title: 代码生成
    details: 从 GraphQL Schema 自动生成 Go 模型、Resolver、DataLoader，减少样板代码
  - icon: 🔌
    title: 开箱即用
    details: 内置数据库连接池、Redis、队列、消息系统、文件存储等常用功能
  - icon: 📊
    title: 生产就绪
    details: 健康检查、Prometheus 指标、优雅关闭，满足生产环境需求
  - icon: 🛡️
    title: 类型安全
    details: 基于 gqlgen，享受 Go 的类型安全和高性能
---

## 技术栈

| 组件 | 技术 |
|------|------|
| GraphQL | gqlgen |
| ORM | GORM (MySQL) |
| 缓存 | Redis |
| 消息队列 | asynq (Redis-based) |
| 消息中间件 | NATS JetStream |
| 存储 | S3/MinIO, 腾讯云 COS |
| 日志 | zerolog |
| 认证 | JWT |
| 数据库迁移 | Atlas |

## 快速体验

```bash
# 安装 CLI
go install github.com/light-speak/lighthouse@latest

# 创建项目
lighthouse generate:init --module github.com/myorg/myapp --models user,post

# 进入项目
cd myapp

# 启动服务
go run . app:start
```

访问 [http://localhost:8080](http://localhost:8080) 查看 GraphQL Playground。
