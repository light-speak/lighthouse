# 项目结构

初始化后的项目结构：

```
myproject/
├── schema/              # GraphQL schema 定义
│   ├── schema.graphql   # 基础类型和指令定义
│   ├── user.graphql     # 用户相关 schema
│   └── *.graphql        # 其他业务 schema
├── graph/               # gqlgen 生成的代码（勿手动修改）
│   ├── *.generated.go
│   └── federation.go
├── models/              # 数据模型
│   ├── models_gen.go    # 生成的 Go 结构体
│   └── dataloader_gen.go # DataLoader 实现
├── resolver/            # 业务逻辑（手动编写）
│   ├── resolver.go      # Resolver 根结构
│   └── *.resolvers.go   # 各业务 resolver
├── commands/            # CLI 命令
│   ├── command.go       # 命令注册
│   ├── app-start.go     # 启动命令
│   ├── migration.go     # 迁移命令
│   └── schema.go        # Schema 导出命令
├── configs/             # 应用配置
│   └── config.go
├── server/              # HTTP 服务配置
│   └── server.go
├── loader/              # GORM Schema 加载器（用于迁移）
│   └── main.go
├── migrations/          # 数据库迁移文件
├── main.go              # 入口
├── go.mod
├── gqlgen.yml           # gqlgen 配置
├── atlas.hcl            # Atlas 迁移配置
└── .env                 # 环境变量
```

## 目录说明

### schema/

GraphQL schema 定义文件，每个业务域一个 `.graphql` 文件。

```graphql
# schema/user.graphql
type User @loader(keys: ["id"]) {
  id: ID!
  name: String!
  posts: [Post!]! @goField(forceResolver: true)
}

extend type Query {
  me: User! @auth
}
```

### graph/

gqlgen 自动生成的代码，**不要手动修改**。

### models/

生成的 Go 结构体和 DataLoader。

### resolver/

业务逻辑实现，手动编写。gqlgen 会生成 resolver 接口，你需要实现具体逻辑。

### commands/

CLI 命令实现，使用 Lighthouse 的命令框架。

### server/

HTTP 服务配置，包括中间件、GraphQL handler 等。

### loader/

用于 Atlas 数据库迁移的 GORM schema 加载器。需要手动添加模型：

```go
var migrateModels = []interface{}{
    &models.User{},
    &models.Post{},
    // 添加所有需要迁移的模型...
}
```

## 文件命名规范

| 类型 | 命名规范 | 示例 |
|------|----------|------|
| Schema | 小写，业务域命名 | `user.graphql`, `order.graphql` |
| Resolver | 业务域 + `.resolvers.go` | `user.resolvers.go` |
| 模型 | 由 schema 自动生成 | `models_gen.go` |
