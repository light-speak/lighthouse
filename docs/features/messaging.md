# 消息系统

Lighthouse 基于 NATS JetStream 提供消息发布/订阅功能。

## 环境变量

```bash
MESSAGING_DRIVER=nats
MESSAGING_URL=localhost:4222
```

## 发布消息

### 原始消息

```go
import "github.com/light-speak/lighthouse/messaging"

broker := messaging.GetBroker()
err := broker.Publish("user.created", []byte(`{"userId": 123}`))
```

### 类型化消息

```go
type UserCreatedEvent struct {
    UserID uint   `json:"userId"`
    Name   string `json:"name"`
}

err := messaging.PublishTyped("user.created", UserCreatedEvent{
    UserID: 123,
    Name:   "John",
})
```

## 订阅模式

| 模式 | 说明 |
|------|------|
| `ModeQueue` | 队列模式：多个消费者负载均衡，每条消息只被一个消费者处理 |
| `ModeBroadcast` | 广播模式：每个消费者都收到所有消息 |

## 订阅消息

### 队列模式（负载均衡）

```go
unsubscribe, err := messaging.GetBroker().Subscribe(
    ctx,
    "user.created",
    func(msg []byte) error {
        // 处理消息
        var event UserCreatedEvent
        json.Unmarshal(msg, &event)
        log.Printf("User created: %d", event.UserID)
        return nil
    },
    messaging.SubscriberOption{
        Mode:          messaging.ModeQueue,
        DurablePrefix: "user-service",
    },
)
```

### 广播模式

```go
unsubscribe, err := messaging.GetBroker().Subscribe(
    ctx,
    "user.created",
    func(msg []byte) error {
        // 每个实例都收到消息
        return nil
    },
    messaging.SubscriberOption{
        Mode: messaging.ModeBroadcast,
    },
)
```

### 类型化订阅

```go
err := messaging.SubscribeTyped(ctx, "user.created", func(event UserCreatedEvent) error {
    log.Printf("User created: %d - %s", event.UserID, event.Name)
    return nil
})
```

## 取消订阅

```go
unsubscribe, _ := broker.Subscribe(ctx, topic, handler)

// 需要时取消订阅
unsubscribe()
```

## 使用场景

### 事件驱动架构

```go
// 用户服务：发布事件
func (r *mutationResolver) Register(ctx context.Context, input RegisterInput) (*User, error) {
    user := createUser(input)

    // 发布用户创建事件
    messaging.PublishTyped("user.created", UserCreatedEvent{
        UserID: user.ID,
        Name:   user.Name,
    })

    return user, nil
}

// 通知服务：订阅事件
func init() {
    messaging.SubscribeTyped(context.Background(), "user.created", func(event UserCreatedEvent) error {
        sendWelcomeEmail(event.UserID)
        return nil
    })
}
```

### 实时通知

```go
// 发送实时通知
func notifyUser(userId uint, message string) {
    messaging.PublishTyped(fmt.Sprintf("user.%d.notification", userId), NotificationEvent{
        Message: message,
    })
}

// 在 Subscription 中订阅
func (r *subscriptionResolver) Notifications(ctx context.Context) (<-chan *Notification, error) {
    userId := auth.GetCtxUserId(ctx)
    ch := make(chan *Notification, 10)

    topic := fmt.Sprintf("user.%d.notification", userId)
    unsubscribe, _ := messaging.GetBroker().Subscribe(ctx, topic, func(msg []byte) error {
        var event NotificationEvent
        json.Unmarshal(msg, &event)
        ch <- &Notification{Message: event.Message}
        return nil
    }, messaging.SubscriberOption{Mode: messaging.ModeBroadcast})

    go func() {
        <-ctx.Done()
        unsubscribe()
        close(ch)
    }()

    return ch, nil
}
```

## 优雅关闭

```go
func (c *Start) OnExit() func() {
    return func() {
        messaging.GetBroker().Close()
    }
}
```

## NATS 部署

### Docker

```bash
docker run -p 4222:4222 -p 8222:8222 nats:latest -js
```

### Docker Compose

```yaml
services:
  nats:
    image: nats:latest
    command: -js
    ports:
      - "4222:4222"
      - "8222:8222"  # 监控端口
```
