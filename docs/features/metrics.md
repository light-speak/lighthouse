# 监控与指标

Lighthouse 内置 Prometheus 指标收集。

## 内置指标

| 指标名 | 类型 | 标签 | 说明 |
|--------|------|------|------|
| `lighthouse_graphql_resolver_duration_seconds` | Histogram | object, field | Resolver 执行耗时 |
| `lighthouse_graphql_operations_total` | Counter | operation, type | GraphQL 操作计数 |

## 初始化指标

```go
// commands/app-start.go
import "github.com/light-speak/lighthouse/metrics"

func (c *Start) Run() error {
    metrics.Init()  // 注册 Prometheus 指标
    server.StartService()
    return nil
}
```

## 使用 MetricsExtension

```go
// server/server.go
import "github.com/light-speak/lighthouse/extensions"

srv := handler.New(graph.NewExecutableSchema(cfg))
srv.Use(extensions.MetricsExtension{})  // 自动收集 resolver 指标
```

## 暴露 Metrics 端点

```go
import "github.com/prometheus/client_golang/prometheus/promhttp"

router.Handle("/metrics", promhttp.Handler())
```

## Prometheus 配置

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'lighthouse'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

## Grafana Dashboard

### 请求延迟

```txt
# 平均延迟
rate(lighthouse_graphql_resolver_duration_seconds_sum[5m])
/ rate(lighthouse_graphql_resolver_duration_seconds_count[5m])

# P99 延迟
histogram_quantile(0.99, rate(lighthouse_graphql_resolver_duration_seconds_bucket[5m]))
```

### 请求量

```txt
# QPS
rate(lighthouse_graphql_operations_total[5m])

# 按操作类型分组
sum by (type) (rate(lighthouse_graphql_operations_total[5m]))
```

### 慢查询

```txt
# 执行时间超过 1 秒的 resolver
histogram_quantile(0.99, rate(lighthouse_graphql_resolver_duration_seconds_bucket[5m])) > 1
```

## 自定义指标

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    UserLoginTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Namespace: "myapp",
            Name:      "user_login_total",
            Help:      "Total user login count",
        },
        []string{"method"},  // 登录方式：password, oauth, etc.
    )

    ActiveUsers = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Namespace: "myapp",
            Name:      "active_users",
            Help:      "Current active users count",
        },
    )

    RequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Namespace: "myapp",
            Name:      "request_duration_seconds",
            Help:      "Request duration in seconds",
            Buckets:   []float64{0.01, 0.05, 0.1, 0.5, 1, 5},
        },
        []string{"endpoint"},
    )
)

func init() {
    prometheus.MustRegister(
        UserLoginTotal,
        ActiveUsers,
        RequestDuration,
    )
}

// 使用
UserLoginTotal.WithLabelValues("password").Inc()
ActiveUsers.Set(float64(count))
RequestDuration.WithLabelValues("/api/users").Observe(duration.Seconds())
```

## 告警规则

```yaml
# prometheus/rules.yml
groups:
  - name: lighthouse
    rules:
      - alert: HighLatency
        expr: histogram_quantile(0.99, rate(lighthouse_graphql_resolver_duration_seconds_bucket[5m])) > 2
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High GraphQL latency"
          description: "P99 latency is above 2 seconds"

      - alert: HighErrorRate
        expr: rate(lighthouse_graphql_operations_total{status="error"}[5m]) > 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate"
          description: "Error rate is above 10%"
```

## Docker Compose 监控栈

```yaml
services:
  app:
    build: .
    ports:
      - "8080:8080"

  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    environment:
      GF_SECURITY_ADMIN_PASSWORD: admin
```
