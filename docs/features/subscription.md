# 实时推送 (Subscription)

GraphQL Subscription 用于实现服务端到客户端的实时推送。

## 定义 Subscription

```graphql
# schema/chat.graphql
extend type Subscription {
  chatMessage: ChatMessage! @auth
  liveCommentAndDanmu(liveId: ID!): LiveCommentAndDanmu! @auth
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

    // 添加订阅者
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

## 订阅管理器

```go
// resolver/subscription_manager.go
type SubscriptionManager struct {
    mu          sync.RWMutex
    subscribers map[uint]chan *models.ChatMessage
}

func NewSubscriptionManager() *SubscriptionManager {
    return &SubscriptionManager{
        subscribers: make(map[uint]chan *models.ChatMessage),
    }
}

func (m *SubscriptionManager) addSubscriber(userId uint, ch chan *models.ChatMessage) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.subscribers[userId] = ch
}

func (m *SubscriptionManager) removeSubscriber(userId uint) {
    m.mu.Lock()
    defer m.mu.Unlock()
    delete(m.subscribers, userId)
}

func (m *SubscriptionManager) Broadcast(userId uint, msg *models.ChatMessage) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    if ch, ok := m.subscribers[userId]; ok {
        select {
        case ch <- msg:
        default:
            // 通道满，丢弃消息
        }
    }
}
```

## 在 Resolver 中使用

```go
// resolver/resolver.go
type Resolver struct {
    LDB                 *databases.LightDatabase
    UserChatMessageChan *SubscriptionManager
}

func NewResolver() *Resolver {
    return &Resolver{
        LDB:                 databases.LightDatabaseClient,
        UserChatMessageChan: NewSubscriptionManager(),
    }
}
```

## 发送消息

```go
func (r *mutationResolver) SendMessage(ctx context.Context, input SendMessageInput) (*models.ChatMessage, error) {
    // 保存消息到数据库
    msg := &models.ChatMessage{
        FromUserID: auth.GetCtxUserId(ctx),
        ToUserID:   input.ToUserID,
        Content:    input.Content,
    }
    db.Create(msg)

    // 推送给接收者
    r.UserChatMessageChan.Broadcast(input.ToUserID, msg)

    return msg, nil
}
```

## WebSocket 配置

```go
// server/server.go
srv := handler.New(graph.NewExecutableSchema(cfg))

srv.AddTransport(transport.Websocket{
    InitFunc:              auth.WebSocketInitFunc,
    KeepAlivePingInterval: 10 * time.Second,
})
```

## 客户端连接

### graphql-ws

```javascript
import { createClient } from 'graphql-ws';

const client = createClient({
  url: 'ws://localhost:8080/query',
  connectionParams: {
    Authorization: `Bearer ${token}`,
  },
});

// 订阅
const unsubscribe = client.subscribe(
  {
    query: `
      subscription {
        chatMessage {
          id
          content
          fromUser {
            name
          }
        }
      }
    `,
  },
  {
    next: (data) => {
      console.log('New message:', data);
    },
    error: (err) => {
      console.error('Subscription error:', err);
    },
    complete: () => {
      console.log('Subscription complete');
    },
  }
);

// 取消订阅
unsubscribe();
```

### Apollo Client

```javascript
import { ApolloClient, InMemoryCache, split, HttpLink } from '@apollo/client';
import { GraphQLWsLink } from '@apollo/client/link/subscriptions';
import { createClient } from 'graphql-ws';
import { getMainDefinition } from '@apollo/client/utilities';

const httpLink = new HttpLink({
  uri: 'http://localhost:8080/query',
});

const wsLink = new GraphQLWsLink(
  createClient({
    url: 'ws://localhost:8080/query',
    connectionParams: {
      Authorization: `Bearer ${token}`,
    },
  })
);

const splitLink = split(
  ({ query }) => {
    const definition = getMainDefinition(query);
    return (
      definition.kind === 'OperationDefinition' &&
      definition.operation === 'subscription'
    );
  },
  wsLink,
  httpLink
);

const client = new ApolloClient({
  link: splitLink,
  cache: new InMemoryCache(),
});
```

## 结合消息系统

使用 NATS 实现跨实例的 Subscription：

```go
func (r *subscriptionResolver) ChatMessage(ctx context.Context) (<-chan *models.ChatMessage, error) {
    userId := auth.GetCtxUserId(ctx)
    ch := make(chan *models.ChatMessage, 10)

    // 订阅 NATS 消息
    topic := fmt.Sprintf("user.%d.chat", userId)
    unsubscribe, _ := messaging.GetBroker().Subscribe(ctx, topic, func(msg []byte) error {
        var chatMsg models.ChatMessage
        json.Unmarshal(msg, &chatMsg)
        ch <- &chatMsg
        return nil
    }, messaging.SubscriberOption{Mode: messaging.ModeBroadcast})

    go func() {
        <-ctx.Done()
        unsubscribe()
        close(ch)
    }()

    return ch, nil
}

// 发送消息时推送到 NATS
func (r *mutationResolver) SendMessage(ctx context.Context, input SendMessageInput) (*models.ChatMessage, error) {
    msg := createMessage(input)

    // 推送到 NATS，所有实例的订阅者都能收到
    topic := fmt.Sprintf("user.%d.chat", input.ToUserID)
    messaging.PublishTyped(topic, msg)

    return msg, nil
}
```
