# 健康检查

Lighthouse 提供 Kubernetes 风格的健康检查端点。

## 端点

| 端点 | 用途 | 检查内容 |
|------|------|----------|
| `/health` | Liveness（存活检查） | 进程是否存活 |
| `/ready` | Readiness（就绪检查） | 数据库、内存、连接池 |

## 配置路由

```go
import "github.com/light-speak/lighthouse/routers"

routers.Config.HeartbeatPath = "/health"
routers.Config.ReadinessPath = "/ready"
```

## 响应格式

### Liveness

```json
{
  "status": "alive",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Readiness

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

## 状态值

| 状态 | HTTP 码 | 说明 |
|------|---------|------|
| `healthy` | 200 | 所有检查通过 |
| `degraded` | 200 | 部分检查警告但仍可用 |
| `unhealthy` | 503 | 关键检查失败 |

## 自定义阈值

```go
import "github.com/light-speak/lighthouse/routers/health"

health.SetConfig(&health.Config{
    // 连接池使用率超过 80% 触发降级
    DBMaxOpenConnsThreshold: 0.8,

    // 内存超过 2GB 触发降级
    MemoryThresholdMB: 2048,

    // 数据库 ping 超时
    DBPingTimeout: 5 * time.Second,
})
```

## Kubernetes 配置

```yaml
apiVersion: v1
kind: Pod
spec:
  containers:
  - name: app
    image: myapp:latest
    ports:
    - containerPort: 8080

    livenessProbe:
      httpGet:
        path: /health
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 10
      timeoutSeconds: 3
      failureThreshold: 3

    readinessProbe:
      httpGet:
        path: /ready
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 10
      timeoutSeconds: 5
      failureThreshold: 3
```

## 检查项说明

### 数据库检查

- 验证数据库连接是否初始化
- 执行 ping 测试连接可用性
- 记录响应延迟

### 内存检查

- 读取当前内存使用量
- 与阈值比较

### 连接池检查

- 获取连接池统计信息
- 计算使用率（in_use / max）
- 超过阈值触发降级

## 自定义检查

可以扩展健康检查添加自定义检查项：

```go
// 添加 Redis 检查
func checkRedis() health.CheckResult {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    if err := redis.Client.Ping(ctx).Err(); err != nil {
        return health.CheckResult{
            Status:  "unhealthy",
            Message: err.Error(),
        }
    }

    return health.CheckResult{
        Status: "healthy",
    }
}
```
