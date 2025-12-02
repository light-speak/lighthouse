# Health Check 健康检查

Lighthouse 提供两个健康检查端点，用于 Kubernetes、负载均衡器等进行服务探测。

## 端点说明

| 端点 | 类型 | 用途 | 检查内容 |
|------|------|------|---------|
| `/health` | Liveness | 存活检查 | 进程是否存活 |
| `/ready` | Readiness | 就绪检查 | 数据库、内存、连接池 |

## 环境变量配置

```bash
MID_HEARTBEAT_PATH=/health    # 存活检查路径
MID_READINESS_PATH=/ready     # 就绪检查路径
```

## 响应格式

### `/health` - 存活检查

```json
{
  "status": "alive",
  "timestamp": "2024-12-02T15:00:00Z"
}
```

- 始终返回 `200 OK`（只要进程存活）

### `/ready` - 就绪检查

```json
{
  "status": "healthy",
  "timestamp": "2024-12-02T15:00:00Z",
  "checks": {
    "database": {
      "status": "healthy",
      "latency": "2.5ms"
    },
    "memory": {
      "status": "healthy",
      "message": "256.0 MB"
    },
    "db_pool": {
      "status": "healthy",
      "message": "in_use=10 idle=40 max=200"
    }
  }
}
```

**状态说明：**

| status | HTTP 状态码 | 含义 |
|--------|------------|------|
| `healthy` | 200 | 所有检查通过，可接收流量 |
| `degraded` | 200 | 部分检查警告（内存/连接池），仍可接收流量 |
| `unhealthy` | 503 | 数据库不可用，应停止接收流量 |

**检查项说明：**

| 检查项 | 检查内容 | 失败条件 |
|--------|---------|---------|
| `database` | 数据库 Ping | 连接失败或超时 (3s) |
| `memory` | 内存使用量 | 超过 1GB (可配置) |
| `db_pool` | 连接池使用率 | 超过 80% (可配置) |

## 自定义配置

```go
import "github.com/light-speak/lighthouse/routers/health"

func init() {
    health.SetConfig(&health.Config{
        DBMaxOpenConnsThreshold: 0.8,           // 连接池使用率阈值 (80%)
        MemoryThresholdMB:       2048,          // 内存阈值 (2GB)
        DBPingTimeout:           5 * time.Second, // 数据库 ping 超时
    })
}
```

## Kubernetes 配置

### 基础配置

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  template:
    spec:
      containers:
      - name: my-app
        image: my-app:latest
        ports:
        - containerPort: 8080

        # 存活探针 - 检查进程是否存活
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10    # 启动后等待时间
          periodSeconds: 10          # 检查间隔
          timeoutSeconds: 3          # 超时时间
          failureThreshold: 3        # 失败次数后重启

        # 就绪探针 - 检查是否可以接收流量
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 5
          failureThreshold: 3        # 失败次数后停止流量
          successThreshold: 1        # 成功次数后恢复流量
```

### 完整示例 (含资源限制)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: live-api
  labels:
    app: live-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: live-api
  template:
    metadata:
      labels:
        app: live-api
    spec:
      containers:
      - name: live-api
        image: live-api:v1.0.0
        ports:
        - containerPort: 8080
          name: http

        env:
        - name: APP_ENV
          value: "production"
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: host

        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "1Gi"
            cpu: "1000m"

        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
          timeoutSeconds: 3
          failureThreshold: 3

        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 5
          failureThreshold: 3

---
apiVersion: v1
kind: Service
metadata:
  name: live-api
spec:
  selector:
    app: live-api
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
```

## Docker Compose 配置

```yaml
version: '3.8'
services:
  api:
    image: my-app:latest
    ports:
      - "8080:8080"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/ready"]
      interval: 10s
      timeout: 5s
      retries: 3
      start_period: 10s
    depends_on:
      db:
        condition: service_healthy

  db:
    image: mysql:8.0
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 5s
      timeout: 3s
      retries: 5
```

## Nginx 负载均衡配置

```nginx
upstream backend {
    server 10.0.0.1:8080 max_fails=3 fail_timeout=30s;
    server 10.0.0.2:8080 max_fails=3 fail_timeout=30s;
    server 10.0.0.3:8080 max_fails=3 fail_timeout=30s;
}

server {
    location / {
        proxy_pass http://backend;

        # 健康检查 (需要 nginx_upstream_check_module)
        health_check interval=5s fails=3 passes=2 uri=/ready;
    }
}
```

## 最佳实践

1. **Liveness vs Readiness**
   - Liveness 失败 → K8s 重启 Pod
   - Readiness 失败 → K8s 停止发送流量，不重启
   - 数据库临时不可用时，应该是 Readiness 失败，而非 Liveness

2. **超时设置**
   - Liveness 超时应短（3s），避免僵尸进程
   - Readiness 超时可稍长（5s），给数据库 ping 足够时间

3. **启动延迟**
   - `initialDelaySeconds` 要大于应用启动时间
   - 避免启动过程中被误判为不健康

4. **降级处理**
   - `degraded` 状态返回 200，仍接收流量
   - 可配合监控告警，提前发现问题
