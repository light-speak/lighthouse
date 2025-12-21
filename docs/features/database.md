# 数据库

Lighthouse 基于 GORM 提供数据库连接池管理，支持主从分离。

## 环境变量配置

```bash
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=myapp
DB_LOG_LEVEL=info

# 连接池配置
DB_MAX_IDLE_CONNS=50      # 最大空闲连接数
DB_MAX_OPEN_CONNS=200     # 最大打开连接数
DB_CONN_MAX_LIFETIME=30   # 连接最大生命周期（分钟）
DB_CONN_MAX_IDLE_TIME=5   # 空闲连接最大存活时间（分钟）

# 主从模式
DB_ENABLE_SLAVE=false
DB_MAIN_HOST=master
DB_SLAVE_HOST=slave1,slave2
```

## 获取数据库连接

```go
// 获取主库（写操作）
db, err := r.LDB.GetDB(ctx)

// 获取从库（读操作，负载均衡）
db, err := r.LDB.GetSlaveDB(ctx)
```

## 基本查询

```go
func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
    db, err := r.LDB.GetSlaveDB(ctx)
    if err != nil {
        return nil, lighterr.NewDatabaseError("服务器繁忙", err)
    }

    user := &models.User{}
    if err := db.First(user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, lighterr.NewNotFoundError("用户不存在")
        }
        return nil, lighterr.NewDatabaseError("查询失败", err)
    }

    return user, nil
}
```

## 事务处理

```go
func (r *mutationResolver) Transfer(ctx context.Context, fromId, toId string, amount int) (bool, error) {
    db, _ := r.LDB.GetDB(ctx)

    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 扣款
    if err := tx.Model(&models.Wallet{}).
        Where("user_id = ? AND balance >= ?", fromId, amount).
        Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
        tx.Rollback()
        return false, err
    }

    // 入账
    if err := tx.Model(&models.Wallet{}).
        Where("user_id = ?", toId).
        Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
        tx.Rollback()
        return false, err
    }

    if err := tx.Commit().Error; err != nil {
        return false, err
    }

    return true, nil
}
```

## 分页查询

```go
func (r *queryResolver) UserList(ctx context.Context, page *int, pageSize *int) ([]*models.User, error) {
    db, _ := r.LDB.GetSlaveDB(ctx)

    query := db.Order("created_at DESC")

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

## 监控连接池

```go
sqlDB, _ := db.DB()
stats := sqlDB.Stats()
logs.Info().
    Int("in_use", stats.InUse).
    Int("idle", stats.Idle).
    Int("open", stats.OpenConnections).
    Int("max", stats.MaxOpenConnections).
    Msg("database pool stats")
```

## 连接池调优

| 参数 | 说明 | 建议值 |
|------|------|--------|
| `DB_MAX_IDLE_CONNS` | 空闲连接数 | 10-50 |
| `DB_MAX_OPEN_CONNS` | 最大连接数 | 100-200 |
| `DB_CONN_MAX_LIFETIME` | 连接生命周期 | 30 分钟 |
| `DB_CONN_MAX_IDLE_TIME` | 空闲超时 | 3-5 分钟 |

::: warning 注意
`DB_CONN_MAX_IDLE_TIME` 是防止连接泄漏的关键参数，必须设置。
:::

## 优雅关闭

在应用退出时关闭数据库连接：

```go
func (c *Start) OnExit() func() {
    return func() {
        logs.Info().Msg("shutting down gracefully...")
        databases.CloseConnections()
    }
}
```
