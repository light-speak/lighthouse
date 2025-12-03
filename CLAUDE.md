# Lighthouse Framework

Lighthouse 是一个基于 Go 的 GraphQL 框架，提供数据库连接池管理、Redis、消息队列等功能。

## 项目结构

```
lighthouse/
├── databases/       # 数据库连接管理 (GORM + MySQL)
├── redis/           # Redis 客户端管理
├── queue/           # 异步任务队列 (asynq)
├── messaging/       # 消息传递 (NATS)
├── routers/         # GraphQL 路由和中间件
│   ├── auth/        # 认证指令
│   ├── dataloader/  # DataLoader 批量查询
│   └── health/      # 健康检查端点
├── lightcmd/        # CLI 和代码生成
│   ├── generate/    # 代码生成器
│   └── initization/ # 项目初始化模板
├── logs/            # 日志模块
├── storages/        # 存储 (S3/COS)
├── templates/       # 模板引擎
├── utils/           # 工具函数
└── lighterr/        # 错误处理
```

## 关键配置

### 数据库连接池

环境变量配置（在 `databases/env.go` 中加载）：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `DB_MAX_IDLE_CONNS` | 10 | 最大空闲连接数 |
| `DB_MAX_OPEN_CONNS` | 100 | 最大打开连接数 |
| `DB_CONN_MAX_LIFETIME` | 30 | 连接最大生命周期（分钟） |
| `DB_CONN_MAX_IDLE_TIME` | 3 | 空闲连接最大存活时间（分钟） |
| `DB_PREPARE_STMT` | false | 是否启用 prepared statement 缓存 |

**重要**：
- `DB_CONN_MAX_IDLE_TIME` 是防止连接泄漏的关键参数
- `DB_PREPARE_STMT=true` 会导致连接累积，仅在高频重复查询场景启用

### Redis 连接池

环境变量配置（在 `redis/env.go` 中加载）：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `REDIS_ENABLE` | false | 是否启用 Redis |
| `REDIS_HOST` | localhost | Redis 主机 |
| `REDIS_PORT` | 6379 | Redis 端口 |
| `REDIS_PASSWORD` | | Redis 密码 |
| `REDIS_DB` | 0 | Redis 数据库 |
| `REDIS_POOL_SIZE` | 10 | 连接池大小 |
| `REDIS_MIN_IDLE_CONNS` | 5 | 最小空闲连接数 |

## 代码生成

### DataLoader 模板

`lightcmd/generate/tpl/dataloader.tpl` 生成的 DataLoader 使用独立 context：

```go
// 所有 fetch 函数使用独立的 context，不受客户端断开影响
queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

这避免了 WebSocket 连接断开时的 "context canceled" 错误。

### 应用启动模板

`lightcmd/initization/tpl/appstart.tpl` 包含 `OnExit()` 实现，确保优雅关闭：

```go
func (c *Start) OnExit() func() {
    return func() {
        logs.Info().Msg("shutting down gracefully...")
        // 关闭数据库、Redis、队列等连接
    }
}
```

## 常见问题

### 连接数持续增长

如果看到 `in_use` 和 `idle` 持续增长：

1. **检查 `PrepareStmt`**：启用时每个连接会缓存 prepared statements
2. **检查 `DB_CONN_MAX_IDLE_TIME`**：必须设置，否则空闲连接不会被清理
3. **确保 `OnExit()` 正确实现**：应用退出时需要调用 `CloseConnections()`

### Context Canceled 错误

在 GraphQL subscriptions 场景下：
- DataLoader 的 fetch 函数使用 `context.Background()` 而非请求 context
- 这样客户端断开不会中断正在执行的数据库查询

### 队列任务并发安全

`queue/queue.go` 使用 `sync.Mutex` 保护 `JobConfigMap`，使用 `sync.Once` 确保客户端单例初始化。

## 开发指南

### 修改模板后重新生成

```bash
# 重新生成 dataloader
light gen dataloader

# 重新生成所有
light gen all
```

### 添加新的数据库

1. 在 `databases/env.go` 配置数据源
2. 使用 `initDatabaseWithRetry()` 初始化
3. 确保在 `OnExit()` 中调用 `CloseConnections()`

## 监控数据库连接

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
