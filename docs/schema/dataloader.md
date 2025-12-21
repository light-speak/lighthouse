# DataLoader

DataLoader 用于解决 GraphQL N+1 查询问题，通过批量加载优化数据库查询。

## 启用 DataLoader

在类型上添加 `@loader` 指令：

```graphql
type User @loader(keys: ["id"]) {
  id: ID!
  name: String!
}
```

## 生成的加载器类型

### 主键加载器

```go
// 单个加载
user, err := models.GetUserIdLoader(ctx).Load(ctx, userId)

// 批量加载
users, err := models.GetUserIdLoader(ctx).LoadAll(ctx, userIds)
```

### 外键加载器（一对一）

```graphql
type Wallet @loader(keys: ["id"], extraKeys: ["userId"]) {
  id: ID!
  userId: ID! @index @unique
  balance: Int!
}
```

```go
// 按 userId 加载单个钱包
wallet, err := models.GetWalletUserIdLoader(ctx).Load(ctx, userId)
```

### 外键列表加载器（一对多）

```graphql
type Post @loader(keys: ["id"], extraKeys: ["userId"]) {
  id: ID!
  userId: ID! @index
  title: String!
}
```

```go
// 按 userId 加载文章列表
posts, err := models.GetPostUserIdListLoader(ctx).Load(ctx, userId)
```

### 复合键加载器

```graphql
type ChatUser @loader(keys: ["chatId", "userId"]) {
  id: ID!
  chatId: ID! @index
  userId: ID! @index
}
```

```go
// 按 chatId + userId 加载
chatUser, err := models.GetChatUserChatIdWithUserIdLoader(ctx).Load(ctx, chatId, userId)
```

## 在 Resolver 中使用

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

### 关联字段示例

```go
// 一对一关系
func (r *userResolver) Wallet(ctx context.Context, obj *models.User) (*models.Wallet, error) {
    return models.GetWalletUserIdLoader(ctx).Load(ctx, obj.ID)
}

// 一对多关系
func (r *userResolver) Posts(ctx context.Context, obj *models.User) ([]*models.Post, error) {
    return models.GetPostUserIdListLoader(ctx).Load(ctx, obj.ID)
}

// 复合键关系
func (r *liveResolver) IsMuted(ctx context.Context, obj *models.Live) (bool, error) {
    userId := auth.GetCtxUserId(ctx)
    if userId == 0 {
        return false, nil
    }
    chatUser, _ := models.GetChatUserChatIdWithUserIdLoader(ctx).Load(ctx, obj.ChatID, userId)
    return chatUser != nil && chatUser.IsMuted, nil
}
```

## DataLoader 中间件

DataLoader 需要中间件支持，在 server.go 中配置：

```go
import "github.com/light-speak/lighthouse/routers/dataloader"

router := routers.NewRouter()
router.Use(dataloader.Middleware(db))
```

## 工作原理

1. **请求收集**：在一个 GraphQL 请求中，多个 resolver 调用同一个 DataLoader
2. **批量合并**：DataLoader 等待一个微小的时间窗口（默认 1ms），收集所有请求
3. **批量查询**：将收集到的 ID 合并成一次批量查询
4. **结果分发**：将查询结果分发给各个等待的 resolver

```
请求: { users { id name wallet { balance } } }

没有 DataLoader:
  SELECT * FROM users WHERE id = 1
  SELECT * FROM wallets WHERE user_id = 1
  SELECT * FROM users WHERE id = 2
  SELECT * FROM wallets WHERE user_id = 2
  ... (2N 次查询)

使用 DataLoader:
  SELECT * FROM users WHERE id IN (1, 2, 3, ...)
  SELECT * FROM wallets WHERE user_id IN (1, 2, 3, ...)
  ... (2 次查询)
```

## 自定义加载逻辑

如果需要自定义加载逻辑，可以在生成的 DataLoader 基础上扩展：

```go
// models/dataloader_custom.go
func GetActiveUserIdLoader(ctx context.Context) *UserLoader {
    // 自定义加载逻辑，只加载活跃用户
    // ...
}
```
