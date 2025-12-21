# Resolver 编写

## Resolver 结构

```go
// resolver/resolver.go
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

## Query Resolver

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

## Mutation Resolver

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

## 字段级 Resolver

### 关联字段 - 使用 DataLoader

```go
func (r *userResolver) Wallet(ctx context.Context, obj *models.User) (*models.Wallet, error) {
    return models.GetWalletUserIdLoader(ctx).Load(ctx, obj.ID)
}

func (r *userResolver) Posts(ctx context.Context, obj *models.User) ([]*models.Post, error) {
    return models.GetPostUserIdListLoader(ctx).Load(ctx, obj.ID)
}
```

### URL 转换

```go
func (r *userResolver) Avatar(ctx context.Context, obj *models.User) (*string, error) {
    if obj.Avatar == nil {
        return nil, nil
    }
    url := r.Storage.GetPublicURL(*obj.Avatar)
    return &url, nil
}
```

### 虚拟字段计算

```go
func (r *userResolver) StarCount(ctx context.Context, obj *models.User) (int, error) {
    db, err := r.LDB.GetSlaveDB(ctx)
    if err != nil {
        return 0, lighterr.NewDatabaseError("服务器繁忙", err)
    }

    count := int64(0)
    db.Model(&models.UserStarLive{}).Where("user_id = ?", obj.ID).Count(&count)
    return int(count), nil
}
```

### 条件计算（当前用户相关）

```go
func (r *liveResolver) Star(ctx context.Context, obj *models.Live) (bool, error) {
    userId := auth.GetCtxUserId(ctx)
    if userId == 0 {
        return false, nil
    }

    userStarLive, _ := models.GetUserStarLiveLiveIdWithUserIdLoader(ctx).Load(ctx, obj.ID, userId)
    return userStarLive != nil, nil
}
```

## Subscription Resolver

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

## 错误处理

使用 `lighterr` 包返回结构化错误：

```go
import "github.com/light-speak/lighthouse/lighterr"

lighterr.NewDatabaseError("操作失败", err)
lighterr.NewBadRequestError("参数错误")
lighterr.NewNotFoundError("数据不存在")
lighterr.NewUnauthorizedError("请先登录")
lighterr.NewForbiddenError("无权访问")
lighterr.NewOperationFailedError("操作失败", err)
```

## 日志记录

```go
import "github.com/light-speak/lighthouse/logs"

logs.Info().Msgf("用户登录: %d", userId)
logs.Error().Err(err).Msg("数据库查询失败")
logs.Debug().Interface("data", obj).Msg("调试信息")
```
