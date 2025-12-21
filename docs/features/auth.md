# 中间件与认证

Lighthouse 提供多种认证中间件，适用于不同场景。

## 中间件类型

| 中间件 | 用途 | 读取的 Header |
|--------|------|---------------|
| `auth.Middleware()` | 用户 JWT 认证 | `Authorization: Bearer <token>` |
| `auth.AdminAuthMiddleware()` | 管理后台认证 | `X-Session-Id`, `RemoteAddr`, `User-Agent` |
| `auth.XUserMiddleware()` | 微服务内部调用 | `X-User-Id` |

## 配置中间件

```go
// server/server.go
package server

import (
    "github.com/go-chi/chi/v5"
    "github.com/light-speak/lighthouse/routers"
    "github.com/light-speak/lighthouse/routers/auth"
)

func NewGraphqlServer(resolver graph.ResolverRoot) *chi.Mux {
    r := routers.NewRouter()

    // 用户端：JWT 认证
    r.Use(auth.Middleware())

    // 或者管理后台：Session 认证
    // r.Use(auth.AdminAuthMiddleware())

    // 或者微服务内部调用：X-User-Id
    // r.Use(auth.XUserMiddleware())

    // ... GraphQL handler setup
    return r
}
```

## 从 Context 获取认证信息

```go
import "github.com/light-speak/lighthouse/routers/auth"

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
    // 获取用户 ID（JWT 或 X-User-Id）
    userId := auth.GetCtxUserId(ctx)
    if userId == 0 {
        return nil, lighterr.NewUnauthorizedError("请先登录")
    }

    // 获取 Session ID（AdminAuthMiddleware）
    sessionId := auth.GetCtxSession(ctx)

    // 获取客户端 IP
    clientIP := auth.GetCtxClientIP(ctx)

    // 获取 User-Agent
    userAgent := auth.GetCtxUserAgent(ctx)

    // ...
}
```

## 辅助函数

```go
// 检查是否登录
if !auth.IsLogin(ctx) {
    return nil, lighterr.NewUnauthorizedError("请先登录")
}

// 检查是否是当前用户
if !auth.IsCurrentUser(ctx, targetUserId) {
    return nil, lighterr.NewForbiddenError("无权访问")
}
```

## JWT Token

### 生成 Token

```go
import "github.com/light-speak/lighthouse/routers/auth"

// 生成 Token
token, err := auth.GetToken(userId)
if err != nil {
    return "", err
}
```

### 验证 Token

```go
// 验证并获取用户 ID
userId, err := auth.GetUserId(token)
if err != nil {
    return 0, err
}
```

### 配置 JWT

```bash
# .env
JWT_SECRET=your-secret-key
JWT_EXPIRE=72h  # Token 过期时间
```

## WebSocket 认证

WebSocket 连接在 `connectionParams` 中传递认证信息：

```go
// server/server.go
srv.AddTransport(transport.Websocket{
    InitFunc: auth.WebSocketInitFunc,
    KeepAlivePingInterval: 10 * time.Second,
})
```

### 客户端连接

```javascript
const client = createClient({
  url: 'ws://localhost:8080/graphql',
  connectionParams: {
    Authorization: 'Bearer <token>',
    // 或者
    'X-User-Id': '123',
  },
});
```

## @auth 指令

框架内置的认证指令：

```graphql
extend type Query {
  me: User! @auth
  privateData: String! @auth(msg: "请先登录后再查看")
}
```

在 server.go 中自动绑定：

```go
cfg.Directives.Auth = auth.AuthDirective
```

## 自定义指令

实现 `@own` 指令检查资源所有权：

```go
cfg.Directives.Own = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
    userId := auth.GetCtxUserId(ctx)

    // 通过反射获取 obj.UserID
    v := reflect.ValueOf(obj).Elem()
    ownerID := v.FieldByName("UserID").Uint()

    if uint(ownerID) != userId {
        return nil, errors.New("无权访问此资源")
    }

    return next(ctx)
}
```

实现 `@hidden` 指令隐藏字段：

```go
cfg.Directives.Hidden = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
    return nil, nil  // 总是返回 nil
}
```
