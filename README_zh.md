# **🚢 Lighthouse GraphQL Framework**

[English](https://github.com/light-speak/lighthouse/blob/main/README.md) | [中文](https://github.com/light-speak/lighthouse/blob/main/README_zh.md)

[![CI](https://github.com/light-speak/lighthouse/actions/workflows/main.yml/badge.svg)](https://github.com/light-speak/lighthouse/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/light-speak/lighthouse/branch/main/graph/badge.svg)](https://codecov.io/gh/light-speak/lighthouse)
[![Go Report Card](https://goreportcard.com/badge/github.com/light-speak/lighthouse)](https://goreportcard.com/report/github.com/light-speak/lighthouse)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

`Lighthouse` 是一个功能丰富的自研 GraphQL 框架，旨在简化基于微服务架构的 GraphQL 服务开发。框架集成了日志系统（使用 `zeroLog`），支持 Elasticsearch、文件日志模式和 Redis 缓存，同时具有灵活的指令系统和强大的自定义配置能力。框架目前内置对 `gorm` 的支持，未来还将扩展更多的 ORM 选项。

## 特性

- **微服务架构支持**：采用独立的微服务模式，不支持 GraphQL Federation，而是通过自定义的服务注册中心来管理微服务。
- **自定义指令**：支持丰富的自定义指令，可以实现动态查询、过滤、关联等操作。
- **可扩展性**：通过配置文件支持多种 GraphQL 文件结构，灵活满足不同项目需求。
- **ORM 集成**：目前支持 `gorm`，未来会增加对更多 ORM 库的支持。
- **日志与缓存集成**：集成了 `zeroLog` 日志系统，支持 Elasticsearch、文件日志和 Redis 缓存。

## 快速开始

### 安装

1. **使用 `go install` 安装**

   ```bash
   go install github.com/light-speak/lighthouse@latest
   ```

2. **创建新项目**

   ```bash
   lighthouse generate:init
   ```

### 配置文件 (lighthouse.yml)

`lighthouse.yml` 是 `lighthouse` 的核心配置文件，用于指定 GraphQL Schema 路径、文件扩展名、ORM 设置等。以下是示例配置：

```yaml
# lighthouse.yml

schema:
  ext:
    - graphql       # 支持的文件扩展名
    - graphqls
  path:
    - schema        # GraphQL Schema 文件所在路径
  model:
    orm: gorm       # ORM 配置，当前支持 gorm
```

- `schema.ext`：指定 Schema 文件的扩展名，可以是 `.graphql` 或 `.graphqls`。
- `schema.path`：定义 Schema 文件的路径，框架将自动加载该路径下的所有文件。
- `model.orm`：当前支持 `gorm` 作为 ORM 库。

### 目录结构

`example` 项目结构如下：

```plaintext
.
├── cmd                     # CLI 相关代码
│   ├── cmd.go              # 主命令入口
│   ├── migrate
│   │   └── migrate.go      # 数据库迁移逻辑
│   └── start
│       └── start.go        # 启动服务入口
├── models                  # 数据模型相关定义
│   ├── enum.go             # 枚举类型定义
│   ├── input.go            # 输入类型定义
│   ├── interface.go        # 接口定义
│   ├── model.go            # 模型结构
│   └── response.go         # 响应数据结构
├── repo                    # 数据库操作封装
│   └── repo.go
├── resolver                # GraphQL 解析器
│   ├── mutation.go         # Mutation 解析
│   ├── query.go            # Query 解析
│   └── resolver.go         # Resolver 主入口
├── schema                  # GraphQL Schema 文件
│   └── user.graphql        # 示例 Schema 文件
└── service                 # 服务逻辑
    └── service.go
```

### 接下来步骤

在 `schema` 目录下填入自定义的 `schema` 文件，然后执行以下命令生成对应的代码：

lighthouse generate:schema

### 使用指令

在 `lighthouse` 中，你可以在 GraphQL Schema 中使用以下指令：

- **@skip / @include**：条件查询指令，用于动态控制字段是否包含在响应中。
- **@enum**：用于定义枚举类型字段，目前只支持 `int8` 类型。
- **@paginate / @find / @first**：用于分页、查找和获取第一个结果的查询。
- **@in / @eq / @neq / @gt / @gte / @lt / @lte / @like / @notIn**：用于参数过滤的指令，支持各种比较运算符。
- **@belongsTo / @hasMany / @hasOne / @morphTo / @morphToMany**：关系映射指令，用于定义模型之间的关系。
- **@index / @unique**：为字段创建索引或添加唯一约束。
- **@defaultString / @defaultInt**：为字段设置默认值。
- **@tag**：用于标记字段的附加属性。
- **@model**：标记类型为数据库模型。
- **@softDeleteModel**：标记类型为数据库模型，并支持软删除功能。
- **@order**：用于对查询结果进行排序。
- **@cache**：用于缓存查询结果，提高响应速度。

### 示例代码

以下是一个示例查询，在获取用户数据时使用了 `@paginate` 指令进行分页：

```graphql
type Query {
  users: [User] @paginate(scopes: ["active"])
}

type User @model(name: "UserModel") {
  id: ID!
  name: String!
  age: Int
  posts: [Post] @hasMany(relation: "Post", foreignKey: "user_id")
}
```

## 扩展与自定义

`lighthouse` 提供了灵活的扩展接口，你可以：

- **添加自定义指令**：编写自己的指令来扩展框架的功能。
- **支持其他 ORM**：参考 `gorm` 集成方式，添加对其他 ORM 库的支持。

## 开发计划

| 🚀 功能分类    | ✨ 功能描述                                          | 📅 状态  |
| -------------- | -------------------------------------------------- | -------- |
| 🛠️ 自定义指令  | 添加自定义指令的支持                                | ✅ 已完成  |
| 📊 查询指令    | 添加 @find 和 @first 注解以支持查询功能             | ✅ 已完成  |
| 🔍 查询与过滤   | 添加日期范围过滤指令                                | ✅ 已完成  |
|                | 添加字符串匹配指令                                  | ✅ 已完成  |
|                | 添加动态排序指令                                    | ✅ 已完成  |
| 📊 分页指令    | 添加 @paginate 注解以支持分页功能                   | ✅ 已完成  |
| 📜 条件查询指令 | 添加 @skip 和 @include 条件查询指令                | 🚧 进行中  |
| 📚 关系映射指令 | 添加 @morphTo, @morphToMany, @hasOne, @manyToMany  | 🚧 进行中  |
| 🔧 微服务管理   | 添加微服务注册中心                                  | ⏳ 计划中  |
| 💾 缓存集成    | 集成 Redis 作为缓存支持                            | ✅ 已完成  |
| 📝 日志系统    | 集成 zeroLog 日志系统，支持 Elasticsearch 和文件日志 | ✅ 已完成  |
| 🔄 缓存指令    | 添加 @cache 指令来支持缓存查询结果                  | 🚧 进行中  |
| 🔀 排序指令    | 添加 @order 指令来支持对查询结果进行排序            | 🚧 进行中  |
| 🗄️ ORM 支持    | 扩展对其他 ORM 的支持，如 `ent`、`sqlc` 等           | ⏳ 计划中  |
| 📑 文档生成工具 | 自动化生成 GraphQL Schema 文档                     | ⏳ 计划中  |
| 📦 插件支持    | 提供插件系统，支持社区贡献和功能扩展                 | ⏳ 计划中  |
| 🌐 前端工具    | 开发类似 Apollo Studio 的前端，用于查询与测试       | ⏳ 计划中  |
| 📊 性能追踪    | 支持每个字段和每个服务的性能追踪，以优化 GraphQL 查询性能 | ⏳ 计划中  |

