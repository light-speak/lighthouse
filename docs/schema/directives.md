# 指令 Directives

## 认证指令 @auth

框架内置的认证指令，检查用户是否登录：

```graphql
extend type Query {
  me: User! @auth
  privateData: String! @auth(msg: "请先登录后再查看")
}
```

在 server.go 中已自动绑定：

```go
cfg := graph.Config{
    Resolvers: &resolver.Resolver{},
}
cfg.Directives.Auth = auth.AuthDirective
```

## 自定义指令 @own / @hidden

这些指令在 schema 中定义，需要自己实现：

```go
// server/server.go
cfg.Directives.Own = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
    userId := auth.GetCtxUserId(ctx)
    // 检查 obj 的 UserId 是否等于当前用户
    // 实现你的所有权逻辑
    return next(ctx)
}

cfg.Directives.Hidden = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
    // 返回 nil 隐藏字段
    return nil, nil
}
```

## 数据库字段指令

| 指令 | 用途 | 示例 |
|------|------|------|
| `@varchar(length: Int!)` | VARCHAR 长度 | `name: String! @varchar(length: 100)` |
| `@text` | TEXT 类型 | `content: String! @text` |
| `@longtext` | LONGTEXT 类型 | `body: String! @longtext` |
| `@index(name: String)` | 创建索引 | `userId: ID! @index` |
| `@unique` | 唯一约束 | `email: String! @unique` |
| `@default(value: String!)` | 默认值 | `status: Int! @default(value: "0")` |
| `@gorm(value: String!)` | GORM 标签 | `count: Int! @gorm(value: "-")` |

### 示例

```graphql
type Article @loader(keys: ["id"]) {
  id: ID!

  # 字符串长度控制
  title: String! @varchar(length: 200)
  summary: String @text
  content: String! @longtext

  # 索引
  authorId: ID! @index
  categoryId: ID! @index(name: "idx_category")

  # 唯一约束
  slug: String! @unique

  # 默认值
  viewCount: Int! @default(value: "0")
  status: Int! @default(value: "1")

  # 虚拟字段（不存数据库）
  author: User! @gorm(value: "-") @goField(forceResolver: true)
}
```

## DataLoader 指令 @loader

### 单键加载

```graphql
type User @loader(keys: ["id"]) {
  id: ID!
  name: String!
}
```

生成的加载器：

```go
models.GetUserIdLoader(ctx).Load(ctx, userId)
models.GetUserIdLoader(ctx).LoadAll(ctx, userIds)
```

### 复合键加载

```graphql
type ChatUser @loader(keys: ["chatId", "userId"]) {
  id: ID!
  chatId: ID! @index
  userId: ID! @index
}
```

生成的加载器：

```go
models.GetChatUserChatIdWithUserIdLoader(ctx).Load(ctx, chatId, userId)
```

### 额外键加载

```graphql
type ChatUser @loader(keys: ["chatId", "userId"], extraKeys: ["userId"]) {
  id: ID!
  chatId: ID! @index
  userId: ID! @index
}
```

额外生成按 userId 查询的列表加载器：

```go
models.GetChatUserUserIdListLoader(ctx).Load(ctx, userId)
```

## Go 代码生成指令

### @goField

```graphql
type User {
  # 强制生成 resolver（用于关联字段、虚拟字段）
  posts: [Post!]! @goField(forceResolver: true)

  # 自定义 Go 字段名
  userId: ID! @goField(name: "UserID")

  # 区分 null 和未传（用于 Input 类型）
  name: String @goField(omittable: true)
}
```

### @goModel

```graphql
# 映射到已有的 Go 类型
scalar Time @goModel(model: "time.Time")

# 映射到自定义类型
type CustomType @goModel(model: "myapp/models.CustomType") {
  # ...
}
```

### @goTag

```graphql
type User {
  # 添加自定义 Go 标签
  email: String! @goTag(key: "validate", value: "email")
  phone: String @goTag(key: "json", value: "phone,omitempty")
}
```
