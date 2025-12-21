package mcp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RegisterResources 注册所有资源
func RegisterResources(s *Server) {
	// 框架文档
	registerFrameworkDocs(s)

	// 指令文档
	registerDirectiveDocs(s)

	// 配置文档
	registerConfigDocs(s)

	// 示例代码
	registerExamples(s)

	// 项目结构
	registerProjectStructure(s)
}

// registerFrameworkDocs 注册框架文档
func registerFrameworkDocs(s *Server) {
	s.RegisterResource(
		"lighthouse://docs/overview",
		"Lighthouse 框架概述",
		"Lighthouse 框架的整体介绍和核心概念",
		"text/markdown",
		func() (*ResourceContent, error) {
			content := `# Lighthouse Framework

Lighthouse 是一个基于 Go 的 GraphQL 框架，提供数据库连接池管理、Redis、消息队列等功能。

## 核心特性

- **Schema-First 开发模式** - 先定义 GraphQL Schema，自动生成 Go 代码
- **自动代码生成** - Models、Resolvers、DataLoaders 自动生成
- **内置认证指令** - @auth, @own 等开箱即用
- **DataLoader 支持** - 自动解决 N+1 查询问题
- **数据库连接池管理** - 支持主从架构
- **异步任务队列** - 基于 asynq 的任务处理
- **消息传递** - NATS 支持

## 快速开始

### 1. 安装
` + "```bash" + `
go install github.com/light-speak/lighthouse@latest
` + "```" + `

### 2. 初始化项目
` + "```bash" + `
lighthouse generate:init --name github.com/yourname/myproject
cd myproject
` + "```" + `

### 3. 定义 Schema
编辑 ` + "`graph/schema.graphql`" + `:
` + "```graphql" + `
type User @loader {
  id: ID!
  name: String! @varchar(length: 100)
  email: String! @unique
  posts: [Post!]!
}

type Post @loader {
  id: ID!
  title: String!
  content: String! @text
  userId: Int!
  user: User!
}

type Query {
  users: [User!]!
  user(id: ID!): User
}
` + "```" + `

### 4. 生成代码
` + "```bash" + `
lighthouse gen
# 或
go run . gen
` + "```" + `

### 5. 启动服务
` + "```bash" + `
go run . app:start
` + "```" + `

## 项目结构

` + "```" + `
myproject/
├── commands/        # CLI 命令
├── graph/           # GraphQL 相关
│   ├── schema.graphql
│   └── resolver.go
├── models/          # 生成的模型
├── server/          # HTTP 服务器
├── tasks/           # 队列任务
├── main.go
└── gqlgen.yml       # gqlgen 配置
` + "```" + `
`
			return &ResourceContent{
				URI:      "lighthouse://docs/overview",
				MimeType: "text/markdown",
				Text:     content,
			}, nil
		},
	)

	// 注册 CLAUDE.md（如果存在）
	s.RegisterResource(
		"lighthouse://docs/claude-md",
		"项目文档 (CLAUDE.md)",
		"当前项目的 CLAUDE.md 文档",
		"text/markdown",
		func() (*ResourceContent, error) {
			// 尝试读取当前目录的 CLAUDE.md
			content, err := os.ReadFile("CLAUDE.md")
			if err != nil {
				// 尝试从 lighthouse 模块目录读取
				if gopath := os.Getenv("GOPATH"); gopath != "" {
					claudePath := filepath.Join(gopath, "pkg/mod/github.com/light-speak/lighthouse@*/CLAUDE.md")
					matches, _ := filepath.Glob(claudePath)
					if len(matches) > 0 {
						content, err = os.ReadFile(matches[len(matches)-1])
					}
				}
				if err != nil {
					return &ResourceContent{
						URI:      "lighthouse://docs/claude-md",
						MimeType: "text/markdown",
						Text:     "# CLAUDE.md 未找到\n\n当前目录下没有 CLAUDE.md 文件。",
					}, nil
				}
			}
			return &ResourceContent{
				URI:      "lighthouse://docs/claude-md",
				MimeType: "text/markdown",
				Text:     string(content),
			}, nil
		},
	)
}

// registerDirectiveDocs 注册指令文档
func registerDirectiveDocs(s *Server) {
	directives := []struct {
		name string
		desc string
	}{
		{"loader", "DataLoader 批量加载指令"},
		{"auth", "认证保护指令"},
		{"own", "所有权校验指令"},
		{"hidden", "隐藏字段指令"},
		{"varchar", "设置 varchar 长度"},
		{"text", "TEXT 类型字段"},
		{"longtext", "LONGTEXT 类型字段"},
		{"gorm", "GORM tag 指令"},
		{"index", "索引指令"},
		{"unique", "唯一索引指令"},
		{"default", "默认值指令"},
	}

	// 注册所有指令文档
	s.RegisterResource(
		"lighthouse://docs/directives",
		"所有 GraphQL 指令",
		"Lighthouse 支持的所有 GraphQL 指令列表",
		"text/markdown",
		func() (*ResourceContent, error) {
			var content strings.Builder
			content.WriteString("# Lighthouse GraphQL 指令\n\n")
			content.WriteString("| 指令 | 说明 |\n")
			content.WriteString("|------|------|\n")
			for _, d := range directives {
				content.WriteString(fmt.Sprintf("| @%s | %s |\n", d.name, d.desc))
			}
			content.WriteString("\n---\n\n")
			for _, d := range directives {
				content.WriteString(getDirectiveInfo(d.name))
				content.WriteString("\n---\n\n")
			}
			return &ResourceContent{
				URI:      "lighthouse://docs/directives",
				MimeType: "text/markdown",
				Text:     content.String(),
			}, nil
		},
	)

	// 为每个指令注册单独的资源
	for _, d := range directives {
		name := d.name
		desc := d.desc
		s.RegisterResource(
			fmt.Sprintf("lighthouse://docs/directives/%s", name),
			fmt.Sprintf("@%s 指令", name),
			desc,
			"text/markdown",
			func() (*ResourceContent, error) {
				return &ResourceContent{
					URI:      fmt.Sprintf("lighthouse://docs/directives/%s", name),
					MimeType: "text/markdown",
					Text:     getDirectiveInfo(name),
				}, nil
			},
		)
	}
}

// registerConfigDocs 注册配置文档
func registerConfigDocs(s *Server) {
	configs := []struct {
		name string
		desc string
	}{
		{"database", "数据库配置"},
		{"redis", "Redis 配置"},
		{"queue", "队列配置"},
		{"messaging", "消息队列配置"},
		{"cors", "CORS 配置"},
		{"health", "健康检查配置"},
	}

	// 注册所有配置文档
	s.RegisterResource(
		"lighthouse://docs/config",
		"所有配置项",
		"Lighthouse 支持的所有配置项",
		"text/markdown",
		func() (*ResourceContent, error) {
			return &ResourceContent{
				URI:      "lighthouse://docs/config",
				MimeType: "text/markdown",
				Text:     getConfigInfo("all"),
			}, nil
		},
	)

	// 为每个配置模块注册单独的资源
	for _, c := range configs {
		name := c.name
		desc := c.desc
		s.RegisterResource(
			fmt.Sprintf("lighthouse://docs/config/%s", name),
			fmt.Sprintf("%s 模块配置", name),
			desc,
			"text/markdown",
			func() (*ResourceContent, error) {
				return &ResourceContent{
					URI:      fmt.Sprintf("lighthouse://docs/config/%s", name),
					MimeType: "text/markdown",
					Text:     getConfigInfo(name),
				}, nil
			},
		)
	}
}

// registerExamples 注册示例代码
func registerExamples(s *Server) {
	s.RegisterResource(
		"lighthouse://examples/schema",
		"GraphQL Schema 示例",
		"完整的 GraphQL Schema 示例",
		"text/plain",
		func() (*ResourceContent, error) {
			content := `# Lighthouse GraphQL Schema 示例

## 基础类型定义

` + "```graphql" + `
# 用户类型，启用 DataLoader
type User @loader {
  id: ID!
  name: String! @varchar(length: 100)
  email: String! @unique @varchar(length: 255)
  password: String! @hidden
  avatar: String @varchar(length: 500)
  bio: String @text
  status: Int! @default(value: "1")
  createdAt: Time!
  updatedAt: Time!
  deletedAt: DeletedAt

  # 关联
  posts: [Post!]!
  comments: [Comment!]!
}

# 帖子类型
type Post @loader(keys: ["userId"]) {
  id: ID!
  title: String! @varchar(length: 200)
  content: String! @longtext
  userId: Int! @index
  viewCount: Int! @default(value: "0")
  createdAt: Time!
  updatedAt: Time!
  deletedAt: DeletedAt

  # 关联
  user: User!
  comments: [Comment!]!
}

# 评论类型（多态关联示例）
type Comment @loader(morphKey: "commentableType", unionTypes: ["Post", "Video"]) {
  id: ID!
  content: String! @text
  userId: Int! @index
  commentableId: Int!
  commentableType: String! @varchar(length: 50)
  createdAt: Time!

  user: User!
}
` + "```" + `

## 查询定义

` + "```graphql" + `
type Query {
  # 公开查询
  posts(limit: Int, offset: Int): [Post!]!
  post(id: ID!): Post

  # 需要认证的查询
  me: User! @auth
  myPosts: [Post!]! @auth
  users: [User!]! @auth(msg: "需要管理员权限")
}
` + "```" + `

## Mutation 定义

` + "```graphql" + `
type Mutation {
  # 用户认证
  register(input: RegisterInput!): AuthPayload!
  login(email: String!, password: String!): AuthPayload!

  # 帖子操作（需要认证）
  createPost(input: CreatePostInput!): Post! @auth
  updatePost(id: ID!, input: UpdatePostInput!): Post! @auth
  deletePost(id: ID!): Boolean! @auth

  # 评论操作
  createComment(postId: ID!, content: String!): Comment! @auth
}

input RegisterInput {
  name: String!
  email: String!
  password: String!
}

input CreatePostInput {
  title: String!
  content: String!
}

input UpdatePostInput {
  title: String
  content: String
}

type AuthPayload {
  token: String!
  user: User!
}
` + "```" + `

## Subscription 定义

` + "```graphql" + `
type Subscription {
  postCreated: Post!
  commentAdded(postId: ID!): Comment!
}
` + "```" + `
`
			return &ResourceContent{
				URI:      "lighthouse://examples/schema",
				MimeType: "text/plain",
				Text:     content,
			}, nil
		},
	)

	s.RegisterResource(
		"lighthouse://examples/resolver",
		"Resolver 示例",
		"GraphQL Resolver 实现示例",
		"text/plain",
		func() (*ResourceContent, error) {
			content := `# Lighthouse Resolver 示例

## 基础 Resolver

` + "```go" + `
package graph

import (
	"context"
	"myproject/models"
	"github.com/light-speak/lighthouse/databases"
	"github.com/light-speak/lighthouse/routers/dataloader"
)

type Resolver struct{}

// Query resolvers

func (r *queryResolver) Posts(ctx context.Context, limit *int, offset *int) ([]*models.Post, error) {
	var posts []*models.Post
	query := databases.GetDB()

	if limit != nil {
		query = query.Limit(*limit)
	}
	if offset != nil {
		query = query.Offset(*offset)
	}

	if err := query.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *queryResolver) Post(ctx context.Context, id string) (*models.Post, error) {
	var post models.Post
	if err := databases.GetDB().First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *queryResolver) Me(ctx context.Context) (*models.User, error) {
	userId := auth.GetCtxUserId(ctx)
	var user models.User
	if err := databases.GetDB().First(&user, userId).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
` + "```" + `

## 使用 DataLoader

` + "```go" + `
// Post.user 字段使用 DataLoader 避免 N+1
func (r *postResolver) User(ctx context.Context, obj *models.Post) (*models.User, error) {
	loader := dataloader.GetLoaderFromCtx[*models.User](ctx, "UserLoader")
	return loader.Load(ctx, obj.UserID)
}

// User.posts 字段
func (r *userResolver) Posts(ctx context.Context, obj *models.User) ([]*models.Post, error) {
	loader := dataloader.GetLoaderFromCtx[*models.Post](ctx, "PostLoaderByUserID")
	return loader.Load(ctx, obj.ID)
}
` + "```" + `

## Mutation Resolver

` + "```go" + `
func (r *mutationResolver) CreatePost(ctx context.Context, input models.CreatePostInput) (*models.Post, error) {
	userId := auth.GetCtxUserId(ctx)

	post := &models.Post{
		Title:   input.Title,
		Content: input.Content,
		UserID:  userId,
	}

	if err := databases.GetDB().Create(post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func (r *mutationResolver) UpdatePost(ctx context.Context, id string, input models.UpdatePostInput) (*models.Post, error) {
	var post models.Post
	if err := databases.GetDB().First(&post, id).Error; err != nil {
		return nil, err
	}

	// 检查所有权
	userId := auth.GetCtxUserId(ctx)
	if post.UserID != userId {
		return nil, errors.New("无权修改此帖子")
	}

	if input.Title != nil {
		post.Title = *input.Title
	}
	if input.Content != nil {
		post.Content = *input.Content
	}

	if err := databases.GetDB().Save(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}
` + "```" + `
`
			return &ResourceContent{
				URI:      "lighthouse://examples/resolver",
				MimeType: "text/plain",
				Text:     content,
			}, nil
		},
	)

	s.RegisterResource(
		"lighthouse://examples/dataloader",
		"DataLoader 示例",
		"DataLoader 使用示例",
		"text/plain",
		func() (*ResourceContent, error) {
			content := `# Lighthouse DataLoader 示例

## Schema 定义

在 GraphQL schema 中使用 @loader 指令：

` + "```graphql" + `
# 基础用法 - 通过 ID 加载
type User @loader {
  id: ID!
  name: String!
}

# 指定额外的查询键
type Post @loader(keys: ["userId", "categoryId"]) {
  id: ID!
  title: String!
  userId: Int!
  categoryId: Int!
}

# 多态关联
type Comment @loader(morphKey: "commentableType", unionTypes: ["Post", "Video"]) {
  id: ID!
  content: String!
  commentableId: Int!
  commentableType: String!
}
` + "```" + `

## 生成的 DataLoader

运行 ` + "`lighthouse gen`" + ` 后会在 ` + "`models/dataloader_gen.go`" + ` 生成：

` + "```go" + `
// UserLoader - 通过 ID 加载
type UserLoader struct {
	*dataloadgen.Loader[int, *User]
}

func NewUserLoader(db *gorm.DB) *UserLoader {
	return &UserLoader{
		Loader: dataloadgen.NewLoader(func(ctx context.Context, ids []int) ([]*User, []error) {
			// 使用独立的 context，不受客户端断开影响
			queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			var users []*User
			if err := db.WithContext(queryCtx).Where("id IN ?", ids).Find(&users).Error; err != nil {
				// 返回错误
			}
			// 按 ID 排序返回
			return users, nil
		}),
	}
}

// PostLoaderByUserID - 通过 userId 加载
type PostLoaderByUserID struct {
	*dataloadgen.Loader[int, []*Post]
}
` + "```" + `

## 在 Resolver 中使用

` + "```go" + `
package graph

import (
	"context"
	"myproject/models"
	"github.com/light-speak/lighthouse/routers/dataloader"
)

func (r *postResolver) User(ctx context.Context, obj *models.Post) (*models.User, error) {
	// 从 context 获取 loader
	loader := dataloader.GetLoaderFromCtx[*models.User](ctx, "UserLoader")

	// Load 会自动批量查询
	return loader.Load(ctx, obj.UserID)
}

func (r *userResolver) Posts(ctx context.Context, obj *models.User) ([]*models.Post, error) {
	loader := dataloader.GetLoaderFromCtx[[]*models.Post](ctx, "PostLoaderByUserID")
	return loader.Load(ctx, obj.ID)
}
` + "```" + `

## 注册 DataLoader 中间件

在 server 启动时注册：

` + "```go" + `
package server

import (
	"myproject/models"
	"github.com/light-speak/lighthouse/databases"
	"github.com/light-speak/lighthouse/routers/dataloader"
)

func init() {
	// 注册所有 loader
	dataloader.RegisterLoader(&models.UserLoader{})
	dataloader.RegisterLoader(&models.PostLoader{})
	dataloader.RegisterLoader(&models.PostLoaderByUserID{})
}

func StartService() {
	db := databases.GetDB()

	// 添加 DataLoader 中间件
	router.Use(dataloader.Middleware(db))

	// ... 其他配置
}
` + "```" + `

## 最佳实践

1. **始终使用 DataLoader 加载关联** - 避免 N+1 查询
2. **DataLoader 使用独立 context** - 不受客户端断开影响
3. **批量查询结果要按输入顺序返回** - DataLoader 要求
4. **处理缺失数据** - 返回 nil 而不是错误
`
			return &ResourceContent{
				URI:      "lighthouse://examples/dataloader",
				MimeType: "text/plain",
				Text:     content,
			}, nil
		},
	)

	s.RegisterResource(
		"lighthouse://examples/command",
		"CLI 命令示例",
		"自定义 CLI 命令实现示例",
		"text/plain",
		func() (*ResourceContent, error) {
			content := `# Lighthouse CLI 命令示例

## 生成命令

` + "```bash" + `
lighthouse generate:command --name migrate --scope db
` + "```" + `

会在 ` + "`commands/`" + ` 目录下生成：

## 命令实现

` + "```go" + `
package commands

import (
	"github.com/light-speak/lighthouse/lightcmd/cmd"
	"github.com/light-speak/lighthouse/logs"
)

type MigrateCommand struct{}

func (c *MigrateCommand) Name() string {
	return "db:migrate"
}

func (c *MigrateCommand) Usage() string {
	return "Run database migrations"
}

func (c *MigrateCommand) Args() []*cmd.CommandArg {
	return []*cmd.CommandArg{
		{
			Name:     "step",
			Type:     cmd.Int,
			Usage:    "Number of migrations to run",
			Required: false,
			Default:  0,
		},
		{
			Name:     "rollback",
			Type:     cmd.Bool,
			Usage:    "Rollback migrations",
			Required: false,
			Default:  false,
		},
	}
}

func (c *MigrateCommand) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		step := *flagValues["step"].(*int)
		rollback := *flagValues["rollback"].(*bool)

		if rollback {
			logs.Info().Msgf("Rolling back %d migrations...", step)
			// 执行回滚逻辑
		} else {
			logs.Info().Msg("Running migrations...")
			// 执行迁移逻辑
		}

		return nil
	}
}

func (c *MigrateCommand) OnExit() func() {
	return func() {
		logs.Info().Msg("Migration command cleanup...")
	}
}

func init() {
	cmd.AddCommand(&MigrateCommand{})
}
` + "```" + `

## 使用命令

` + "```bash" + `
# 交互模式
lighthouse db:migrate

# 命令行模式
lighthouse db:migrate --step 5

# 回滚
lighthouse db:migrate --rollback --step 1

# 查看帮助
lighthouse db:migrate --help
` + "```" + `
`
			return &ResourceContent{
				URI:      "lighthouse://examples/command",
				MimeType: "text/plain",
				Text:     content,
			}, nil
		},
	)

	s.RegisterResource(
		"lighthouse://examples/task",
		"队列任务示例",
		"异步任务实现示例",
		"text/plain",
		func() (*ResourceContent, error) {
			content := `# Lighthouse 队列任务示例

## 生成任务

` + "```bash" + `
lighthouse generate:task --name send_email
` + "```" + `

## 任务实现

` + "```go" + `
package tasks

import (
	"context"
	"encoding/json"

	"github.com/light-speak/lighthouse/logs"
	"github.com/light-speak/lighthouse/queue"
)

// SendEmailPayload 任务载荷
type SendEmailPayload struct {
	To      string ` + "`json:\"to\"`" + `
	Subject string ` + "`json:\"subject\"`" + `
	Body    string ` + "`json:\"body\"`" + `
}

// SendEmailTask 发送邮件任务
type SendEmailTask struct{}

func (t *SendEmailTask) Name() string {
	return "send_email"
}

func (t *SendEmailTask) Execute(ctx context.Context, payload []byte) error {
	var p SendEmailPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return err
	}

	logs.Info().
		Str("to", p.To).
		Str("subject", p.Subject).
		Msg("Sending email...")

	// 实际发送邮件的逻辑
	// ...

	return nil
}

func init() {
	queue.RegisterJob("send_email", queue.JobConfig{
		Name:     "send_email",
		Priority: 1,
		Executor: &SendEmailTask{},
	})
}
` + "```" + `

## 投递任务

` + "```go" + `
package service

import (
	"encoding/json"
	"myproject/tasks"
	"github.com/light-speak/lighthouse/queue"
)

func SendWelcomeEmail(userEmail string) error {
	payload := tasks.SendEmailPayload{
		To:      userEmail,
		Subject: "Welcome to our platform!",
		Body:    "Thank you for registering...",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return queue.Enqueue("send_email", data)
}
` + "```" + `

## 启动队列 Worker

` + "```go" + `
package commands

import (
	"github.com/light-speak/lighthouse/lightcmd/cmd"
	"github.com/light-speak/lighthouse/queue"
)

type QueueWorker struct{}

func (c *QueueWorker) Name() string {
	return "queue:work"
}

func (c *QueueWorker) Usage() string {
	return "Start the queue worker"
}

func (c *QueueWorker) Args() []*cmd.CommandArg {
	return nil
}

func (c *QueueWorker) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		return queue.StartQueue()
	}
}

func (c *QueueWorker) OnExit() func() {
	return func() {
		queue.CloseClient()
	}
}

func init() {
	cmd.AddCommand(&QueueWorker{})
}
` + "```" + `

## 运行 Worker

` + "```bash" + `
lighthouse queue:work
` + "```" + `
`
			return &ResourceContent{
				URI:      "lighthouse://examples/task",
				MimeType: "text/plain",
				Text:     content,
			}, nil
		},
	)
}

// registerProjectStructure 注册项目结构资源
func registerProjectStructure(s *Server) {
	s.RegisterResource(
		"lighthouse://docs/structure",
		"项目结构",
		"Lighthouse 项目的标准目录结构",
		"text/markdown",
		func() (*ResourceContent, error) {
			content := `# Lighthouse 项目结构

## 标准目录结构

` + "```" + `
myproject/
├── commands/           # CLI 命令
│   ├── start.go       # 启动服务命令
│   ├── migrate.go     # 数据库迁移命令
│   └── ...
│
├── graph/              # GraphQL 相关
│   ├── schema.graphql # GraphQL Schema 定义
│   ├── resolver.go    # Resolver 实现
│   └── generated.go   # gqlgen 生成的代码
│
├── models/             # 数据模型
│   ├── models_gen.go  # gqlgen 生成的模型
│   ├── dataloader_gen.go # 生成的 DataLoader
│   └── custom.go      # 自定义模型方法
│
├── server/             # HTTP 服务器
│   └── server.go      # 服务器配置和启动
│
├── tasks/              # 队列任务
│   ├── send_email.go
│   └── ...
│
├── migrations/         # 数据库迁移文件
│   └── ...
│
├── .env               # 环境变量
├── .env.example       # 环境变量示例
├── gqlgen.yml         # gqlgen 配置
├── go.mod
└── main.go            # 入口文件
` + "```" + `

## 各目录说明

### commands/
存放 CLI 命令，每个命令实现 ` + "`CommandInterface`" + ` 接口。

### graph/
存放 GraphQL Schema 和 Resolver。
- ` + "`schema.graphql`" + ` - 定义类型、查询、变更
- ` + "`resolver.go`" + ` - 实现 Resolver 逻辑

### models/
存放数据模型。
- ` + "`*_gen.go`" + ` - 自动生成，不要手动修改
- 可以添加自定义文件扩展模型

### server/
HTTP 服务器配置。
- 路由配置
- 中间件注册
- GraphQL Handler

### tasks/
异步任务定义，使用 asynq 队列。

## 配置文件

### gqlgen.yml
` + "```yaml" + `
schema:
  - graph/*.graphql

exec:
  filename: graph/generated.go
  package: graph

model:
  filename: models/models_gen.go
  package: models

resolver:
  layout: follow-schema
  dir: graph
  package: graph

autobind:
  - "myproject/models"

models:
  ID:
    model:
      - github.com/99designs/gqlgen/graphql.ID
  Int:
    model:
      - github.com/99designs/gqlgen/graphql.Int
  Time:
    model:
      - github.com/99designs/gqlgen/graphql.Time
  DeletedAt:
    model:
      - github.com/light-speak/lighthouse/lightcmd/scalars.DeletedAt
` + "```" + `

### .env
` + "```bash" + `
# Server
PORT=8080
GIN_MODE=release

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=
DB_NAME=myproject
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100
DB_CONN_MAX_LIFETIME=30
DB_CONN_MAX_IDLE_TIME=3

# Redis
REDIS_ENABLE=true
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# CORS
CORS_ALLOW_ORIGINS=*
` + "```" + `
`
			return &ResourceContent{
				URI:      "lighthouse://docs/structure",
				MimeType: "text/markdown",
				Text:     content,
			}, nil
		},
	)
}
