# 快速开始

## 安装 CLI

```bash
go install github.com/light-speak/lighthouse@latest
```

## 创建项目

```bash
# 创建新项目
lighthouse generate:init --module github.com/myorg/myapp --models user,post,comment
```

初始化会创建完整的项目结构，包括：
- GraphQL schema 文件
- gqlgen 配置
- 命令框架（app:start, migration:apply, schema）
- Atlas 迁移配置
- 环境变量模板

## 配置环境变量

编辑 `.env` 文件：

```bash
# Application
APP_NAME=MyApp
APP_PORT=8080
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=myapp

# Redis (可选)
REDIS_ENABLE=false
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_SECRET=your-secret-key
```

## 启动服务

```bash
cd myapp
go run . app:start
```

访问 [http://localhost:8080](http://localhost:8080) 查看 GraphQL Playground。

## 下一步

- 了解 [CLI 命令](/guide/cli) 的完整用法
- 查看 [项目结构](/guide/project-structure) 详解
- 学习 [GraphQL Schema](/schema/basics) 编写
