# GraphQL Schema 基础

## 基础 Schema

每个项目都需要一个基础 schema 文件定义标量类型和指令：

```graphql
# schema/schema.graphql

# 标量类型
scalar Time
scalar DeletedAt

# ============================================
# Go 代码生成指令
# ============================================

# 模型映射
directive @goModel(
    model: String
    models: [String!]
    forceGenerate: Boolean
) on OBJECT | INPUT_OBJECT | SCALAR | ENUM | INTERFACE | UNION

# 字段控制
directive @goField(
    forceResolver: Boolean  # 强制生成 resolver
    name: String            # Go 字段名
    omittable: Boolean      # 区分 null 和未传
) on INPUT_FIELD_DEFINITION | FIELD_DEFINITION

# Go 标签
directive @goTag(
    key: String!
    value: String
) on INPUT_FIELD_DEFINITION | FIELD_DEFINITION

# ============================================
# 业务指令
# ============================================

# 认证：必须登录才能访问
directive @auth(msg: String) on FIELD_DEFINITION

# 所有权：只能访问自己的数据
directive @own on FIELD_DEFINITION

# 隐藏：响应中不返回
directive @hidden on FIELD_DEFINITION

# ============================================
# 数据库指令
# ============================================

directive @longtext on FIELD_DEFINITION
directive @text on FIELD_DEFINITION
directive @varchar(length: Int!) on FIELD_DEFINITION
directive @index(name: String) on FIELD_DEFINITION
directive @unique on FIELD_DEFINITION
directive @default(value: String!) on FIELD_DEFINITION
directive @gorm(value: String!) on FIELD_DEFINITION

# ============================================
# DataLoader 指令
# ============================================

directive @loader(
    keys: [String!]       # 主键字段
    morphKey: String      # 多态键
    unionTypes: [String!] # 多态类型
    extraKeys: [String!]  # 额外的加载键
) on OBJECT
```

## 业务模型定义

```graphql
# schema/user.graphql

# @loader 启用 DataLoader 批量加载
type User @loader(keys: ["id"]) {
  id: ID!
  createdAt: Time!
  updatedAt: Time!
  deletedAt: DeletedAt

  # 字符串字段
  account: String
  name: String! @varchar(length: 100)
  avatar: String @goField(forceResolver: true)  # URL 转换

  # 数值字段带默认值
  level: Int! @default(value: "1")
  experience: Int! @default(value: "0")

  # 索引字段
  serviceChatSourceId: String @index(name: "service")

  # 关联字段 - 使用 @goField(forceResolver: true)
  wallet: Wallet! @goField(forceResolver: true)
  role: Role! @goField(forceResolver: true)

  # 虚拟字段（不映射到数据库）
  starCount: Int! @gorm(value: "-") @goField(forceResolver: true)
  isOnline: Boolean! @gorm(value: "-") @goField(forceResolver: true)
}

# 枚举类型
enum UserStatus {
  ACTIVE
  BANNED
  DELETED
}
```

## Query 和 Mutation

使用 `extend` 扩展根类型：

```graphql
# Query 扩展
extend type Query {
  me: User! @auth
  userList(page: Int, pageSize: Int): [User!]!
  user(id: ID!): User
}

# Mutation 扩展
extend type Mutation {
  register(input: RegisterInput!): String!
  login(input: LoginInput!): String!
  updateUserInfo(input: UpdateUserInput!): User! @auth
}

# Subscription 扩展
extend type Subscription {
  online: String! @auth
}
```

## Input 类型

```graphql
input RegisterInput {
  account: String!
  password: String!
  name: String!
}

input LoginInput {
  account: String!
  password: String!
}

# Partial Update - 使用 omittable 区分 null 和未传
input UpdateUserInput {
  name: String @goField(omittable: true)
  avatar: String @goField(omittable: true)
}
```

## @goField(forceResolver: true) 使用场景

```graphql
type Live @loader(keys: ["id"]) {
  id: ID!
  userId: ID! @index

  # 1. 关联字段 - 需要从其他表加载
  user: User! @goField(forceResolver: true)
  liveStat: LiveStat @goField(forceResolver: true)

  # 2. URL 转换 - 需要处理存储路径
  cover: String! @goField(forceResolver: true)

  # 3. 虚拟字段 - 动态计算，不存数据库
  star: Boolean! @gorm(value: "-") @goField(forceResolver: true)
  viewerCount: Int! @gorm(value: "-") @goField(forceResolver: true)
  shareUrl: String! @gorm(value: "-") @goField(forceResolver: true)

  # 4. 条件相关 - 依赖当前用户
  isMuted: Boolean! @gorm(value: "-") @goField(forceResolver: true)
}
```

## @goField(omittable: true) 使用场景

区分「传了 null」和「没传」的场景：

```graphql
input UpdateUserInput {
  name: String @goField(omittable: true)
  email: String @goField(omittable: true)
  avatar: String @goField(omittable: true)
}
```

生成的 Go 代码：

```go
type UpdateUserInput struct {
    Name   graphql.Omittable[*string]
    Email  graphql.Omittable[*string]
    Avatar graphql.Omittable[*string]
}

// 使用方式
func (r *mutationResolver) UpdateUser(ctx context.Context, input UpdateUserInput) (*User, error) {
    if input.Name.IsSet() {
        if input.Name.Value() == nil {
            // 用户明确传了 null，清空字段
            user.Name = ""
        } else {
            // 用户传了具体值
            user.Name = *input.Name.Value()
        }
    }
    // !input.Name.IsSet() 表示用户没传，不修改
}
```

## 复合键 DataLoader

```graphql
# 需要按 (chatId, userId) 组合查询
type ChatUser @loader(keys: ["chatId", "userId"], extraKeys: ["userId"]) {
  id: ID!
  chatId: ID! @index
  userId: ID! @index
  isMuted: Boolean! @default(value: "false")
}
```
