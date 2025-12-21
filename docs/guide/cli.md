# CLI 命令

## 项目初始化

```bash
lighthouse generate:init --module <module> --models <models>
```

| 参数 | 说明 | 示例 |
|------|------|------|
| `--module` | Go module 路径 | `github.com/myorg/myapp` |
| `--models` | 初始模型列表（逗号分隔） | `user,post,comment` |

## 代码生成

```bash
lighthouse generate:schema
```

从 GraphQL schema 生成：
- Go 模型结构体 (`models/models_gen.go`)
- DataLoader (`models/dataloader_gen.go`)
- Resolver 接口 (`graph/*.generated.go`)

生成过程会输出详细日志：
- 加载的配置文件
- 生成的模型数量和名称
- DataLoader 信息
- 执行 go mod tidy 和 gofmt

## 应用启动

```bash
go run . app:start
```

启动 GraphQL HTTP 服务，默认端口 8080。

## 数据库迁移

```bash
# 生成迁移文件
atlas migrate diff --env dev

# 应用迁移
go run . migration:apply --env=dev
```

## 导出 Schema

```bash
go run . schema
```

导出完整的 GraphQL schema 到 `schema.graphql` 文件。

## 环境变量参考

### 应用配置

```bash
APP_NAME=MyApp
APP_PORT=8080
APP_ENV=development
```

### 数据库配置

```bash
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
```

### Redis 配置

```bash
REDIS_ENABLE=false
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 队列配置

```bash
QUEUE_ENABLE=false
QUEUE_REDIS_HOST=localhost
QUEUE_REDIS_PORT=6379
QUEUE_REDIS_PASSWORD=
QUEUE_REDIS_DB=0
```

### 消息系统配置

```bash
MESSAGING_DRIVER=nats
MESSAGING_URL=localhost:4222
```

### 存储配置

```bash
# S3/MinIO
STORAGE_DRIVER=s3
S3_ENDPOINT=localhost:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_DEFAULT_BUCKET=uploads
S3_USE_SSL=false

# 腾讯云 COS
STORAGE_DRIVER=cos
COS_SECRET_ID=xxx
COS_SECRET_KEY=xxx
COS_BUCKET=bucket-appid
COS_REGION=ap-guangzhou
```
