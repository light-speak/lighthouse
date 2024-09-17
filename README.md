# Lighthouse


---

## 🚀 特色功能

- 🧩 **支持 GraphQL Federation**：方便多个服务整合，轻松扩展功能。
- ⚡ **超快性能**：基于 Go 的高效执行，充分利用多核性能。
- 🔄 **热更新**：在开发环境下无需重启服务器即可自动更新 GraphQL Schema。
- 🎯 **强大的 Dataloader 支持**：彻底解决 N+1 查询问题，性能优化轻松搞定。

---

## 快速开始

本文以Level服务为例，为已有的User附加Level功能，已有UserSchema如下：

```graphql
type User implements BaseModelSoftDelete @key(fields: "id") {
    id: ID!
    name: String!
    phone: String!

    createdAt: Time!
    updatedAt: Time!
    deletedAt: Time
}

extend type Query {
    me: User @resolve
}

extend type Mutation {
    login(phone: String!): String! @resolve
}

```

- 创建并进入项目目录：`mkdir level && cd level`
- 创建 `tools.go` 文件

```go
//go:build tools
// +build tools

package tools

import (
	_ "github.com/light-speak/lighthouse"
)

```

- 进行 `go.mod` 初始化

```shell
go mod init gitlab.staticoft.com/cos/level
```

- 若存在未上线的Go模块，则使用本地replace 例如:
- 此时我为本地开发环境，正式使用不加这行

```
replace github.com/light-speak/lighthouse => ../lighthouse
```

- 执行 `go mod tidy`
- 执行 `go run github.com/light-speak/lighthouse init `

会有一部分文件自动创建，目录如下：

```shell
├── .env
├── .gitignore
├── go.mod
├── go.sum
├── storage
│   └── logs
│       └── log-2024-08-18.log
└── tools.go

```

env文件包含数据库、服务端口、日志等多项设置，优先级低于环境变量

- 执行 `go run github.com/light-speak/lighthouse gql:init `

创建了GraphQL基础文件，目录如下：

```graphql
├── .env
├── .gitignore
├── go.mod
├── go.sum
├── gqlgen.yml
├── graph
│   ├── lighthouse.graphqls
│   ├── resolver.go
│   └── server.go
├── storage
│   └── logs
│       └── log-2024-08-18.log
└── tools.go

```

可直接在这步开始编写GraphQL Schema

- 在graph目录下创建level.graphqls

> SchemaType 需要实现两种类型接口： 1、BaseModel 2、BaseModelSoftDelete

```graphql

type Level implements BaseModel @key(fields: "id") {
    id: ID!
    createdAt: Time!
    updatedAt: Time!

    "等级名称"
    name: String!
}

type UserLevel implements BaseModel @key(fields: "id") {
    id: ID!
    createdAt: Time!
    updatedAt: Time!

    levelId: ID!
    "等级"
    level: Level! @requires(fields: "levelId")

    userId: ID!
    "用户"
    user: User! @requires(fields: "userId")  # 由UserLevel提供UserID字段，交由User服务处理
}

type User @key(fields: "id") @extends {
    id: ID! @external
    level: UserLevel! @provides(fields: "id")  # 由User提供Id字段，Level服务处理
}

extend type Query {
    levels(page: Int! @page, size: Int! @size): [Level!]! @all
}

```

- 执行 `go run github.com/light-speak/lighthouse gql:generate `
- 创建 `main.go` 文件

```go
package main

import (
	"gitlab.staticoft.com/cos/level/graph"
	"github.com/light-speak/lighthouse/log"
)

func main() {

	log.Error("%s", graph.StartServer())
}

```

此时，基础的Resolver已经生成完毕，需要处理Entity.resolvers.go的逻辑

已知生成的代编写逻辑代码如下：

```go
// FindLevelByID is the resolver for the findLevelByID field.
func (r *entityResolver) FindLevelByID(ctx context.Context, id int64) (*Level, error) {
panic(fmt.Errorf("not implemented: FindLevelByID - findLevelByID"))
}

// FindUserByID is the resolver for the findUserByID field.
func (r *entityResolver) FindUserByID(ctx context.Context, id int64) (*User, error) {
panic(fmt.Errorf("not implemented: FindUserByID - findUserByID"))
}

// FindUserLevelByID is the resolver for the findUserLevelByID field.
func (r *entityResolver) FindUserLevelByID(ctx context.Context, id int64) (*UserLevel, error) {
panic(fmt.Errorf("not implemented: FindUserLevelByID - findUserLevelByID"))
}
```

### 逐一解析

#### FindLevelByID

该方法为外界访问暴露端口，可能存在高频读取，例如GraphQL查询中的N+1问题，此时需要使用Dataloader解决，Dataloader也已经自动生成完毕，直接调用：

```go
func (r *entityResolver) FindLevelByID(ctx context.Context, id int64) (*Level, error) {
level, err := For(ctx).FindLevelById.Load(id)
if err != nil {
return nil, err
}
return r.mergeLevel(ctx, level)
}
```

由于GraphQL存在循环调用依赖，为了满足依赖覆盖，故所有返回单个Model的，使用r.merge{{ Model }}(model) （`r.mergeLevel(level)`）
进行依赖融合

返回列表的情况使用 r.merge{{ Model }}s(model) （`r.mergeLevels(level)`）


#### FindUserByID

该方法为Level服务提供根据UserId提供UserLevel模型，对应的是，User -> UserLevel ，所以说人话就是，使用此处的UserId查询出对应UserLevel，并且附加到*User模型，进行返回

```go
func (r *entityResolver) FindUserByID(ctx context.Context, id int64) (*User, error) {
	userLevel, err := r.getUserLevelByUserId(ctx, id)
	if err != nil {
		return nil, err
	}
	return &User{ID: id, UserLevel: userLevel}, nil
}
```

如代码所示，在Level服务，为User提供了UserLevel模型

#### FindUserLevelByID

该方法为查找UserLevel下的User,同理，附加User

```go
func (r *entityResolver) FindUserLevelByID(ctx context.Context, id int64) (*UserLevel, error) {
	userLevel, err := For(ctx).FindUserLevelById.Load(id)
	if err != nil {
		return nil, err
	}
	userLevel.User = &User{ID: userLevel.UserID}
	return r.mergeUserLevel(ctx, userLevel)
}
```

- `go run .` 开始运行吧

## Apollo Router

> 待补充

## Directive

注解实现进度

- [x] @all
- [x] @first   ---- 还未实现Scope
- [x] @eq
```graphql
extend type Query {
    post(title: String! @eq): Post @first
}
```
- [x] @create
- [x] @update
```graphql
extend type Mutation {
    post(title: String!): Post! @create
    editPost(id: ID!, title: String!): Post! @update
}
```

- [x] @page  ----- 还未实现Field内注解
- [x] @size
```graphql
extend type Query {
    posts(page: Int! @page size: Int! @size): [Post!]! @all(scopes: ["postCustomScope"])
}
```

- [x] @count
- [x] @sum
```graphql
extend type Query {
    postCount: Int! @count(model: "Post", scopes: [""])
    postLikeSum: Int! @sum(model: "Post", column: "like", scopes: [""])
}
```
- [x] @requires 将某个作为提供者
```graphql
type UserLevel implements BaseModel @key(fields: "id") {
    id: ID!
    createdAt: Time!
    updatedAt: Time!

    levelId: ID!
    "等级"
    level: Level! @requires(fields: "levelId")

    userId: ID!
    "用户"
    user: User! @requires(fields: "userId") 
}
```
- [x] @resolve 自定义操作，不会被Gen覆盖

- [ ] @auth  接口要求用户登录
- [ ] @userId 注入登录UserId
- [ ] @cache 根据时间缓存该接口
- [ ] @neq 不等于
- [ ] @like 如名
- [ ] @inject 从Context注入参数
- [ ] @validator 正则检验

## N+1问题？

gen已生成Dataloader代码，只需：

```go
posts, err := For(ctx).FindPostById.Load(id)
```

既可解决传统N+1问题
## 其他功能
### RPC通信
> 待补充
### 队列
> 待补充