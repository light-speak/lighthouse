# Lighthouse Framework Documentation

本文档详细介绍 Lighthouse GraphQL 框架的使用方法，供 AI 辅助开发参考。

## 目录

1. [框架概览](#框架概览)
2. [CLI 命令](#cli-命令)
3. [项目结构](#项目结构)
4. [GraphQL Schema 编写](#graphql-schema-编写)
5. [Directives 指令](#directives-指令)
6. [Resolver 编写](#resolver-编写)
7. [DataLoader 使用](#dataloader-使用)
8. [数据库操作](#数据库操作)
9. [数据库迁移](#数据库迁移)
10. [中间件与认证](#中间件与认证)
11. [健康检查](#健康检查)
12. [异步任务队列](#异步任务队列)
13. [消息系统](#消息系统)
14. [文件存储](#文件存储)
15. [实时推送 (Subscription)](#实时推送-subscription)
16. [监控与指标](#监控与指标)
17. [最佳实践](#最佳实践)

---

## 框架概览

### 技术栈

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

### 框架模块

```
lighthouse/
├── lightcmd/      # CLI 工具与代码生成（核心）
├── databases/     # 数据库连接管理（主从支持）
├── redis/         # Redis 缓存
├── queue/         # 异步任务队列
├── messaging/     # NATS 消息系统
├── routers/       # HTTP 路由和中间件
├── lighterr/      # 统一错误处理
├── logs/          # 日志系统
├── storages/      # 文件存储
├── templates/     # 代码模板引擎
└── utils/         # 工具函数
```

---

## CLI 命令

### 项目初始化

```bash
# 创建新项目
lighthouse generate:init --module github.com/myorg/myproject --models user,post,comment
```

初始化会创建完整的项目结构，包括：
- GraphQL schema 文件
- gqlgen 配置
- 命令框架（app:start, migration:apply, schema）
- Atlas 迁移配置
- 环境变量模板

### 代码生成

```bash
# 生成 GraphQL schema、models、resolver、dataloader
lighthouse generate:schema
```

生成过程会输出详细日志：
- 加载的配置文件
- 生成的模型数量和名称
- DataLoader 信息
- 执行 go mod tidy 和 gofmt

### 服务启动

```bash
# 启动 GraphQL 服务
go run . app:start
```

### 数据库迁移

```bash
# 生成迁移文件
atlas migrate diff --env dev

# 应用迁移
go run . migration:apply --env=dev
```

### 导出 Schema

```bash
# 导出完整的 GraphQL schema 到 schema.graphql
go run . schema
```

---

## 项目结构

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

---

## GraphQL Schema 编写

### 基础 Schema (schema/schema.graphql)

```graphql
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

# 字段控制（重要！）
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

### 业务模型定义

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

# 输入类型 - 使用 omittable 区分 null 和未传
input UpdateUserInput {
  name: String @goField(omittable: true)
  avatar: String @goField(omittable: true)
}

# Query 扩展
extend type Query {
  me: User! @auth
  userList(page: Int, pageSize: Int): [User!]!
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

### @goField(forceResolver: true) 使用场景

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

### @goField(omittable: true) 使用场景

```graphql
# Partial Update 场景
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

### 复合键 DataLoader

```graphql
# 需要按 (chatId, userId) 组合查询
type ChatUser @loader(keys: ["chatId", "userId"], extraKeys: ["userId"]) {
  id: ID!
  chatId: ID! @index
  userId: ID! @index
  isMuted: Boolean! @default(value: "false")
}
```

---

## Directives 指令

### 认证指令 @auth

```graphql
extend type Query {
  me: User! @auth
  privateData: String! @auth(msg: "请先登录后再查看")
}
```

框架已内置 `@auth` 指令实现，在 server.go 中绑定：

```go
cfg := graph.Config{
    Resolvers: &resolver.Resolver{},
}
cfg.Directives.Auth = auth.AuthDirective  // 绑定指令
```

### 自定义指令 @own / @hidden

`@own` 和 `@hidden` 在 schema 中定义但需要自己实现：

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

### 数据库字段指令

| 指令 | 用途 | 示例 |
|------|------|------|
| `@varchar(length: Int!)` | VARCHAR 长度 | `name: String! @varchar(length: 100)` |
| `@text` | TEXT 类型 | `content: String! @text` |
| `@longtext` | LONGTEXT 类型 | `body: String! @longtext` |
| `@index(name: String)` | 创建索引 | `userId: ID! @index` |
| `@unique` | 唯一约束 | `email: String! @unique` |
| `@default(value: String!)` | 默认值 | `status: Int! @default(value: "0")` |
| `@gorm(value: String!)` | GORM 标签 | `count: Int! @gorm(value: "-")` |

### DataLoader 指令 @loader

```graphql
# 单键加载
type User @loader(keys: ["id"]) { ... }

# 复合键加载
type ChatUser @loader(keys: ["chatId", "userId"]) { ... }

# 额外键（生成按该键查询的加载器）
type ChatUser @loader(keys: ["chatId", "userId"], extraKeys: ["userId"]) { ... }
```

---

## Resolver 编写

### Resolver 结构 (resolver/resolver.go)

```go
package resolver

import (
    "github.com/light-speak/lighthouse/databases"
    "github.com/light-speak/lighthouse/storages"
)

type Resolver struct {
    LDB     *databases.LightDatabase  // 数据库连接池
    Storage storages.Storage          // 文件存储
}
```

### Query Resolver

```go
func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
    userId := auth.GetCtxUserId(ctx)

    // 使用从库查询
    db, err := r.LDB.GetSlaveDB(ctx)
    if err != nil {
        return nil, lighterr.NewDatabaseError("服务器繁忙", err)
    }

    user := &models.User{}
    if err := db.First(user, userId).Error; err != nil {
        return nil, lighterr.NewNotFoundError("用户不存在")
    }

    return user, nil
}

func (r *queryResolver) UserList(ctx context.Context, page *int, pageSize *int) ([]*models.User, error) {
    db, err := r.LDB.GetSlaveDB(ctx)
    if err != nil {
        return nil, lighterr.NewDatabaseError("服务器繁忙", err)
    }

    query := db.Order("created_at DESC")

    // 分页
    if page != nil && pageSize != nil {
        offset := (*page - 1) * (*pageSize)
        query = query.Offset(offset).Limit(*pageSize)
    }

    users := make([]*models.User, 0)
    if err := query.Find(&users).Error; err != nil {
        return nil, lighterr.NewDatabaseError("查询失败", err)
    }

    return users, nil
}
```

### Mutation Resolver

```go
func (r *mutationResolver) Register(ctx context.Context, input models.RegisterInput) (string, error) {
    db, err := r.LDB.GetDB(ctx)  // 使用主库写入
    if err != nil {
        return "", lighterr.NewDatabaseError("服务器繁忙", err)
    }

    // 开启事务
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 检查账号是否存在
    existingUser := &models.User{}
    if err := tx.Where("account = ?", input.Account).First(existingUser).Error; err == nil {
        tx.Rollback()
        return "", lighterr.NewBadRequestError("账号已存在")
    }

    // 创建用户
    user := &models.User{
        Account: &input.Account,
        Name:    input.Name,
    }
    if err := tx.Create(user).Error; err != nil {
        tx.Rollback()
        return "", lighterr.NewDatabaseError("创建用户失败", err)
    }

    // 提交事务
    if err := tx.Commit().Error; err != nil {
        return "", lighterr.NewDatabaseError("注册失败", err)
    }

    // 生成 JWT Token
    token, err := auth.GetToken(user.ID)
    if err != nil {
        return "", lighterr.NewOperationFailedError("生成令牌失败", err)
    }

    return token, nil
}
```

### 字段级 Resolver

```go
// 关联字段 - 使用 DataLoader 避免 N+1
func (r *userResolver) Wallet(ctx context.Context, obj *models.User) (*models.Wallet, error) {
    return models.GetWalletUserIdLoader(ctx).Load(ctx, obj.ID)
}

// URL 转换
func (r *userResolver) Avatar(ctx context.Context, obj *models.User) (*string, error) {
    if obj.Avatar == nil {
        return nil, nil
    }
    url := r.Storage.GetPublicURL(*obj.Avatar)
    return &url, nil
}

// 虚拟字段计算
func (r *userResolver) StarCount(ctx context.Context, obj *models.User) (int, error) {
    db, err := r.LDB.GetSlaveDB(ctx)
    if err != nil {
        return 0, lighterr.NewDatabaseError("服务器繁忙", err)
    }

    count := int64(0)
    db.Model(&models.UserStarLive{}).Where("user_id = ?", obj.ID).Count(&count)
    return int(count), nil
}

// 条件计算（当前用户相关）
func (r *liveResolver) Star(ctx context.Context, obj *models.Live) (bool, error) {
    userId := auth.GetCtxUserId(ctx)
    if userId == 0 {
        return false, nil
    }

    userStarLive, _ := models.GetUserStarLiveLiveIdWithUserIdLoader(ctx).Load(ctx, obj.ID, userId)
    return userStarLive != nil, nil
}
```

---

## DataLoader 使用

### 生成的 DataLoader 类型

```go
// 单键加载器
models.GetUserIdLoader(ctx).Load(ctx, userId)
models.GetUserIdLoader(ctx).LoadAll(ctx, userIds)

// 外键加载器（一对一）
models.GetWalletUserIdLoader(ctx).Load(ctx, userId)

// 外键列表加载器（一对多）
models.GetPostImagePostIdListLoader(ctx).Load(ctx, postId)

// 复合键加载器
models.GetChatUserChatIdWithUserIdLoader(ctx).Load(ctx, chatId, userId)

// 额外键加载器
models.GetChatUserUserIdListLoader(ctx).Load(ctx, userId)
```

### 避免 N+1 查询

```go
// ❌ 错误方式：每次循环都查询数据库
for _, live := range lives {
    user := &models.User{}
    db.First(user, live.UserID)  // N 次查询
}

// ✅ 正确方式：使用 DataLoader 批量加载
func (r *liveResolver) User(ctx context.Context, obj *models.Live) (*models.User, error) {
    // 多个请求会被自动合并成一次批量查询
    return models.GetUserIdLoader(ctx).Load(ctx, obj.UserID)
}
```

---

## 数据库操作

### 获取数据库连接

```go
// 获取主库（写操作）
db, err := r.LDB.GetDB(ctx)

// 获取从库（读操作，负载均衡）
db, err := r.LDB.GetSlaveDB(ctx)
```

### 事务处理

```go
func (r *mutationResolver) ComplexOperation(ctx context.Context) (string, error) {
    db, _ := r.LDB.GetDB(ctx)

    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.Create(&model1).Error; err != nil {
        tx.Rollback()
        return "", err
    }

    if err := tx.Create(&model2).Error; err != nil {
        tx.Rollback()
        return "", err
    }

    if err := tx.Commit().Error; err != nil {
        return "", err
    }

    return "success", nil
}
```

---

## 中间件与认证

### 中间件类型

Lighthouse 提供三种认证中间件：

| 中间件 | 用途 | 读取的 Header |
|--------|------|---------------|
| `auth.Middleware()` | 用户 JWT 认证 | `Authorization: Bearer <token>` |
| `auth.AdminAuthMiddleware()` | 管理后台认证（Session/IP/UA） | `X-Session-Id`, `RemoteAddr`, `User-Agent` |
| `auth.XUserMiddleware()` | 微服务内部调用（信任模式） | `X-User-Id` |

### 中间件配置 (server/server.go)

```go
package server

import (
    "github.com/go-chi/chi/v5"
    "github.com/light-speak/lighthouse/routers"
    "github.com/light-speak/lighthouse/routers/auth"
)

func NewGraphqlServer(resolver graph.ResolverRoot) *chi.Mux {
    r := routers.NewRouter()

    // 用户端：JWT 认证
    r.Use(auth.Middleware())

    // 或者管理后台：Session 认证
    // r.Use(auth.AdminAuthMiddleware())

    // 或者微服务内部调用：X-User-Id
    // r.Use(auth.XUserMiddleware())

    // ... GraphQL handler setup
    return r
}
```

### 从 Context 获取认证信息

```go
import "github.com/light-speak/lighthouse/routers/auth"

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
    // 获取用户 ID（JWT 或 X-User-Id）
    userId := auth.GetCtxUserId(ctx)
    if userId == 0 {
        return nil, lighterr.NewUnauthorizedError("请先登录")
    }

    // 获取 Session ID（AdminAuthMiddleware）
    sessionId := auth.GetCtxSession(ctx)

    // 获取客户端 IP
    clientIP := auth.GetCtxClientIP(ctx)

    // 获取 User-Agent
    userAgent := auth.GetCtxUserAgent(ctx)

    // 检查是否登录
    if !auth.IsLogin(ctx) {
        return nil, lighterr.NewUnauthorizedError("请先登录")
    }

    // 检查是否是当前用户
    if !auth.IsCurrentUser(ctx, targetUserId) {
        return nil, lighterr.NewForbiddenError("无权访问")
    }

    // ...
}
```

### WebSocket 认证

WebSocket 连接在 `connectionParams` 中传递认证信息：

```go
// server/server.go
srv.AddTransport(transport.Websocket{
    InitFunc: auth.WebSocketInitFunc,
    // ...
})
```

客户端连接时：

```javascript
const client = createClient({
  url: 'ws://localhost:8080/graphql',
  connectionParams: {
    Authorization: 'Bearer <token>',
    // 或者
    'X-User-Id': '123',
  },
});
```

### JWT Token 生成

```go
import "github.com/light-speak/lighthouse/routers/auth"

// 生成 Token
token, err := auth.GetToken(userId)

// 验证并获取用户 ID
userId, err := auth.GetUserId(token)
```

---

## 健康检查

### 端点配置

Lighthouse 提供两个健康检查端点：

| 端点 | 用途 | 检查内容 |
|------|------|----------|
| `/health` | Liveness（存活检查） | 进程是否存活 |
| `/ready` | Readiness（就绪检查） | 数据库、内存、连接池 |

### 配置路由 (routers/config.go)

```go
routers.Config.HeartbeatPath = "/health"
routers.Config.ReadinessPath = "/ready"
```

### 健康检查响应

就绪检查返回详细的健康状态：

```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z",
  "checks": {
    "database": {
      "status": "healthy",
      "latency": "1.234ms"
    },
    "memory": {
      "status": "healthy",
      "message": "512.0 MB"
    },
    "db_pool": {
      "status": "healthy",
      "message": "in_use=5 idle=10 max=100"
    }
  }
}
```

状态值：
- `healthy`: 所有检查通过
- `degraded`: 部分检查警告但仍可用
- `unhealthy`: 关键检查失败

### 自定义阈值

```go
import "github.com/light-speak/lighthouse/routers/health"

health.SetConfig(&health.Config{
    DBMaxOpenConnsThreshold: 0.8,       // 连接池使用率 80% 触发降级
    MemoryThresholdMB:       2048,      // 内存超过 2GB 触发降级
    DBPingTimeout:           5 * time.Second,
})
```

### Kubernetes 探针配置

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
```

---

## 异步任务队列

基于 [asynq](https://github.com/hibiken/asynq) 的 Redis 异步任务队列。

### 环境变量

```bash
QUEUE_ENABLE=true
QUEUE_REDIS_HOST=localhost
QUEUE_REDIS_PORT=6379
QUEUE_REDIS_PASSWORD=
QUEUE_REDIS_DB=0
```

### 定义任务

```go
package jobs

import (
    "context"
    "encoding/json"
    "github.com/hibiken/asynq"
    "github.com/light-speak/lighthouse/queue"
)

const TypeEmailDelivery = "email:delivery"

type EmailPayload struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

type EmailJob struct{}

func (j *EmailJob) Execute(ctx context.Context, task *asynq.Task) error {
    var payload EmailPayload
    if err := json.Unmarshal(task.Payload(), &payload); err != nil {
        return err
    }

    // 发送邮件逻辑
    return sendEmail(payload.To, payload.Subject, payload.Body)
}

func init() {
    queue.RegisterJob(TypeEmailDelivery, queue.JobConfig{
        Name:     TypeEmailDelivery,
        Priority: 10,
        Executor: &EmailJob{},
    })
}
```

### 发送任务

```go
func SendEmailAsync(to, subject, body string) error {
    client, err := queue.GetClient()
    if err != nil {
        return err
    }

    payload, _ := json.Marshal(EmailPayload{
        To:      to,
        Subject: subject,
        Body:    body,
    })

    task := asynq.NewTask(TypeEmailDelivery, payload)

    // 立即执行
    _, err = client.Enqueue(task)

    // 延迟执行
    _, err = client.Enqueue(task, asynq.ProcessIn(10*time.Minute))

    // 定时执行
    _, err = client.Enqueue(task, asynq.ProcessAt(time.Now().Add(24*time.Hour)))

    return err
}
```

### 启动队列消费者

```go
// commands/queue.go
func (c *QueueWorker) Run() error {
    return queue.StartQueue()
}
```

### 优雅关闭

```go
func (c *Start) OnExit() func() {
    return func() {
        queue.CloseClient()
    }
}
```

---

## 消息系统

基于 NATS JetStream 的消息发布/订阅系统。

### 环境变量

```bash
MESSAGING_DRIVER=nats
MESSAGING_URL=localhost:4222
```

### 发布消息

```go
import "github.com/light-speak/lighthouse/messaging"

// 发布原始消息
broker := messaging.GetBroker()
err := broker.Publish("user.created", []byte(`{"userId": 123}`))

// 发布类型化消息（自动 JSON 序列化）
type UserCreatedEvent struct {
    UserID uint   `json:"userId"`
    Name   string `json:"name"`
}

err := messaging.PublishTyped("user.created", UserCreatedEvent{
    UserID: 123,
    Name:   "John",
})
```

### 订阅消息

两种订阅模式：

| 模式 | 说明 |
|------|------|
| `ModeQueue` | 队列模式：多个消费者负载均衡，每条消息只被一个消费者处理 |
| `ModeBroadcast` | 广播模式：每个消费者都收到所有消息 |

```go
import "github.com/light-speak/lighthouse/messaging"

// 队列模式订阅（负载均衡）
unsubscribe, err := messaging.GetBroker().Subscribe(
    ctx,
    "user.created",
    func(msg []byte) error {
        // 处理消息
        return nil
    },
    messaging.SubscriberOption{
        Mode:          messaging.ModeQueue,
        DurablePrefix: "user-service",
    },
)

// 广播模式订阅
unsubscribe, err := messaging.GetBroker().Subscribe(
    ctx,
    "user.created",
    func(msg []byte) error {
        // 每个实例都收到
        return nil
    },
    messaging.SubscriberOption{
        Mode: messaging.ModeBroadcast,
    },
)

// 类型化订阅
err := messaging.SubscribeTyped(ctx, "user.created", func(event UserCreatedEvent) error {
    log.Printf("User created: %d - %s", event.UserID, event.Name)
    return nil
})

// 取消订阅
unsubscribe()
```

### 优雅关闭

```go
func (c *Start) OnExit() func() {
    return func() {
        messaging.GetBroker().Close()
    }
}
```

---

## 文件存储

统一的文件存储接口，支持 S3/MinIO 和腾讯云 COS。

### 环境变量

```bash
# S3/MinIO
STORAGE_DRIVER=s3
S3_ENDPOINT=localhost:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_DEFAULT_BUCKET=uploads
S3_USE_SSL=false

# 或腾讯云 COS
STORAGE_DRIVER=cos
COS_SECRET_ID=xxx
COS_SECRET_KEY=xxx
COS_BUCKET=bucket-appid
COS_REGION=ap-guangzhou
```

### 初始化存储

```go
import "github.com/light-speak/lighthouse/storages"

storage, err := storages.NewStorage()
```

### 存储接口

```go
type Storage interface {
    // 上传文件
    Put(ctx context.Context, key string, reader io.Reader) error

    // 获取预签名下载 URL
    GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)

    // 获取预签名上传 URL（客户端直传）
    GetPresignedPutURL(ctx context.Context, key string, expiry time.Duration) (string, error)

    // 获取公开访问 URL
    GetPublicURL(key string) string
}
```

### 使用示例

```go
// 在 Resolver 中使用
func (r *Resolver) Storage storages.Storage

// 上传文件
func (r *mutationResolver) UploadAvatar(ctx context.Context, file graphql.Upload) (string, error) {
    key := fmt.Sprintf("avatars/%d/%s", time.Now().Unix(), file.Filename)
    err := r.Storage.Put(ctx, key, file.File)
    if err != nil {
        return "", err
    }
    return key, nil
}

// 字段 Resolver 中转换 URL
func (r *userResolver) Avatar(ctx context.Context, obj *models.User) (*string, error) {
    if obj.Avatar == nil {
        return nil, nil
    }
    url := r.Storage.GetPublicURL(*obj.Avatar)
    return &url, nil
}

// 获取预签名上传 URL（客户端直传）
func (r *mutationResolver) GetUploadURL(ctx context.Context, filename string) (string, error) {
    key := fmt.Sprintf("uploads/%d/%s", time.Now().Unix(), filename)
    return r.Storage.GetPresignedPutURL(ctx, key, 15*time.Minute)
}
```

---

## 数据库迁移

### Atlas 配置 (atlas.hcl)

```hcl
data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./loader",
  ]
}

env "dev" {
  src = data.external_schema.gorm.url
  dev = "mysql://root:@127.0.0.1:3306/test"
  url = "mysql://root:@127.0.0.1:3306/myapp"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "production" {
  src = data.external_schema.gorm.url
  url = env("DATABASE_URL")
  migration {
    dir = "file://migrations"
  }
}
```

### loader/main.go

```go
package main

import (
    "myapp/models"
    "ariga.io/atlas-provider-gorm/gormschema"
    "gorm.io/gorm"
)

var migrateModels = []interface{}{
    &models.User{},
    &models.Wallet{},
    // 添加所有需要迁移的模型...
}

func main() {
    option := gormschema.WithConfig(&gorm.Config{
        DisableForeignKeyConstraintWhenMigrating: true,
        IgnoreRelationshipsWhenMigrating:         true,
    })
    stmts, _ := gormschema.New("mysql", option).Load(migrateModels...)
    io.WriteString(os.Stdout, stmts)
}
```

### 迁移工作流

```bash
# 1. 修改 schema/*.graphql
# 2. 生成 Go 代码
lighthouse generate:schema

# 3. 更新 loader/main.go 添加新模型

# 4. 生成迁移 SQL
atlas migrate diff --env dev

# 5. 检查生成的 migrations/TIMESTAMP.sql

# 6. 应用迁移
go run . migration:apply --env=dev
```

---

## 实时推送 (Subscription)

### 定义 Subscription

```graphql
extend type Subscription {
  chatMessage: ChatMessage! @auth
  liveCommentAndDanmu(liveId: ID!): LiveCommentAndDanmu! @auth
}
```

### Subscription Resolver

```go
func (r *subscriptionResolver) ChatMessage(ctx context.Context) (<-chan *models.ChatMessage, error) {
    userId := auth.GetCtxUserId(ctx)
    if userId == 0 {
        return nil, lighterr.NewUnauthorizedError("请先登录")
    }

    ch := make(chan *models.ChatMessage, 10)

    r.UserChatMessageChan.addSubscriber(userId, ch)

    // 监听客户端断开
    go func() {
        <-ctx.Done()
        r.UserChatMessageChan.removeSubscriber(userId)
        close(ch)
    }()

    return ch, nil
}
```

---

## 监控与指标

### Prometheus 指标

Lighthouse 内置 Prometheus 指标收集：

| 指标名 | 类型 | 标签 | 说明 |
|--------|------|------|------|
| `lighthouse_graphql_resolver_duration_seconds` | Histogram | object, field | Resolver 执行耗时 |
| `lighthouse_graphql_operations_total` | Counter | operation, type | GraphQL 操作计数 |

### 初始化指标

```go
// server/server.go 或 commands/app-start.go
import "github.com/light-speak/lighthouse/metrics"

func main() {
    metrics.Init()  // 注册 Prometheus 指标
    // ...
}
```

### 使用 MetricsExtension

```go
import "github.com/light-speak/lighthouse/extensions"

srv := handler.New(graph.NewExecutableSchema(cfg))
srv.Use(extensions.MetricsExtension{})  // 自动收集 resolver 指标
```

### 暴露 Metrics 端点

```go
import "github.com/prometheus/client_golang/prometheus/promhttp"

router.Handle("/metrics", promhttp.Handler())
```

### 自定义指标

```go
import "github.com/prometheus/client_golang/prometheus"

var MyCounter = prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Namespace: "myapp",
        Name:      "custom_events_total",
        Help:      "Custom events count",
    },
    []string{"event_type"},
)

func init() {
    prometheus.MustRegister(MyCounter)
}

// 使用
MyCounter.WithLabelValues("user_login").Inc()
```

---

## 最佳实践

### 1. Schema 组织

- 每个业务域一个 `.graphql` 文件
- 使用 `extend type Query/Mutation/Subscription` 扩展
- 虚拟字段用 `@gorm(value: "-")` 标记
- 需要 resolver 的字段用 `@goField(forceResolver: true)`

### 2. 错误处理

```go
import "github.com/light-speak/lighthouse/lighterr"

lighterr.NewDatabaseError("操作失败", err)
lighterr.NewBadRequestError("参数错误")
lighterr.NewNotFoundError("数据不存在")
lighterr.NewUnauthorizedError("请先登录")
lighterr.NewOperationFailedError("操作失败", err)
```

### 3. 日志记录

```go
import "github.com/light-speak/lighthouse/logs"

logs.Info().Msgf("用户登录: %d", userId)
logs.Error().Err(err).Msg("数据库查询失败")
logs.Debug().Interface("data", obj).Msg("调试信息")
```

### 4. 认证获取

```go
import "github.com/light-speak/lighthouse/routers/auth"

userId := auth.GetCtxUserId(ctx)
if userId == 0 {
    // 未登录
}
```

### 5. 分页模式

```go
func (r *queryResolver) List(ctx context.Context, page *int, pageSize *int) ([]*models.Item, error) {
    db, _ := r.LDB.GetSlaveDB(ctx)

    query := db.Order("created_at DESC")

    if page != nil && pageSize != nil {
        offset := (*page - 1) * (*pageSize)
        query = query.Offset(offset).Limit(*pageSize)
    }

    items := make([]*models.Item, 0)
    query.Find(&items)

    return items, nil
}
```

---

## 常用命令速查

```bash
# 启动服务
go run . app:start

# 生成代码
lighthouse generate:schema

# 数据库迁移
atlas migrate diff --env dev        # 生成迁移
go run . migration:apply --env=dev  # 应用迁移

# 导出 schema
go run . schema
```

---

## 环境变量参考

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
DB_LOG_LEVEL=info
DB_MAX_IDLE_CONNS=50
DB_MAX_OPEN_CONNS=200
DB_CONN_MAX_LIFETIME=30
DB_CONN_MAX_IDLE_TIME=5

# 主从模式
DB_ENABLE_SLAVE=false
DB_MAIN_HOST=master
DB_SLAVE_HOST=slave1,slave2

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-secret-key

# Storage
STORAGE_DRIVER=s3
S3_ENDPOINT=localhost:9000
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_DEFAULT_BUCKET=default

# Queue
QUEUE_ENABLE=false
QUEUE_REDIS_HOST=localhost
QUEUE_REDIS_PORT=6379

# Messaging
MESSAGING_DRIVER=nats
MESSAGING_URL=localhost:4222
```
