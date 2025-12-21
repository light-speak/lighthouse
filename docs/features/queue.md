# 异步任务队列

Lighthouse 基于 [asynq](https://github.com/hibiken/asynq) 提供 Redis 异步任务队列。

## 环境变量

```bash
QUEUE_ENABLE=true
QUEUE_REDIS_HOST=localhost
QUEUE_REDIS_PORT=6379
QUEUE_REDIS_PASSWORD=
QUEUE_REDIS_DB=0
```

## 定义任务

```go
// jobs/email.go
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

## 注册任务

在 `main.go` 或单独的 `jobs/init.go` 中导入任务包：

```go
import _ "myapp/jobs"
```

## 发送任务

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
    return err
}
```

### 延迟执行

```go
// 10 分钟后执行
_, err = client.Enqueue(task, asynq.ProcessIn(10*time.Minute))
```

### 定时执行

```go
// 指定时间执行
_, err = client.Enqueue(task, asynq.ProcessAt(time.Now().Add(24*time.Hour)))
```

### 设置重试

```go
// 最多重试 3 次
_, err = client.Enqueue(task, asynq.MaxRetry(3))
```

### 设置超时

```go
// 任务超时时间
_, err = client.Enqueue(task, asynq.Timeout(5*time.Minute))
```

## 启动队列消费者

创建队列工作进程命令：

```go
// commands/queue.go
package commands

import (
    "github.com/light-speak/lighthouse/queue"
    _ "myapp/jobs"  // 导入任务
)

type QueueWorker struct{}

func (c *QueueWorker) Name() string {
    return "queue:work"
}

func (c *QueueWorker) Description() string {
    return "Start queue worker"
}

func (c *QueueWorker) Run() error {
    return queue.StartQueue()
}
```

启动：

```bash
go run . queue:work
```

## 任务优先级

优先级数值越高，优先处理：

```go
queue.RegisterJob("critical:task", queue.JobConfig{
    Name:     "critical:task",
    Priority: 100,  // 高优先级
    Executor: &CriticalTaskJob{},
})

queue.RegisterJob("normal:task", queue.JobConfig{
    Name:     "normal:task",
    Priority: 10,   // 普通优先级
    Executor: &NormalTaskJob{},
})
```

## 优雅关闭

```go
func (c *Start) OnExit() func() {
    return func() {
        queue.CloseClient()
    }
}
```

## 监控

asynq 提供 Web UI 监控：

```bash
# 安装 asynqmon
go install github.com/hibiken/asynqmon@latest

# 启动监控
asynqmon --redis-addr=localhost:6379
```

访问 [http://localhost:8080](http://localhost:8080) 查看任务状态。
