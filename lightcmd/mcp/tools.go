package mcp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RegisterTools 注册所有 Lighthouse 工具
func RegisterTools(s *Server) {
	// 代码生成工具
	registerGenerateSchema(s)
	registerGenerateDataloader(s)
	registerGenerateCommand(s)
	registerGenerateTask(s)
	registerInitProject(s)

	// 信息查询工具
	registerGetDirectiveInfo(s)
	registerGetConfigInfo(s)
	registerListGenerators(s)
	registerSearchDocs(s)

	// 项目工具
	registerRunLighthouseCommand(s)
}

// registerGenerateSchema 注册 schema 生成工具
func registerGenerateSchema(s *Server) {
	s.RegisterTool(
		"generate_schema",
		"生成 GraphQL schema 和 models。读取 .graphql 文件并生成对应的 Go 代码，包括 models、resolvers 和 dataloaders。",
		map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
		func(args map[string]interface{}) (*CallToolResult, error) {
			cmd := exec.Command("go", "run", ".", "gen")
			output, err := cmd.CombinedOutput()
			if err != nil {
				return &CallToolResult{
					Content: []ContentBlock{NewTextContent(fmt.Sprintf("生成失败:\n%s\n错误: %s", string(output), err.Error()))},
					IsError: true,
				}, nil
			}
			return &CallToolResult{
				Content: []ContentBlock{NewTextContent(fmt.Sprintf("Schema 生成成功:\n%s", string(output)))},
			}, nil
		},
	)
}

// registerGenerateDataloader 注册 dataloader 生成工具
func registerGenerateDataloader(s *Server) {
	s.RegisterTool(
		"generate_dataloader",
		"生成 DataLoader 代码。DataLoader 用于批量加载数据，避免 N+1 查询问题。需要在 GraphQL schema 中使用 @loader 指令标记模型。",
		map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
		func(args map[string]interface{}) (*CallToolResult, error) {
			cmd := exec.Command("go", "run", ".", "gen")
			output, err := cmd.CombinedOutput()
			if err != nil {
				return &CallToolResult{
					Content: []ContentBlock{NewTextContent(fmt.Sprintf("生成失败:\n%s", string(output)))},
					IsError: true,
				}, nil
			}
			return &CallToolResult{
				Content: []ContentBlock{NewTextContent(fmt.Sprintf("DataLoader 生成成功:\n%s", string(output)))},
			}, nil
		},
	)
}

// registerGenerateCommand 注册命令生成工具
func registerGenerateCommand(s *Server) {
	s.RegisterTool(
		"generate_command",
		"生成新的 CLI 命令。命令会放在 commands/ 目录下，自动注册到命令列表中。",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "命令名称，例如: migrate, seed",
				},
				"scope": map[string]interface{}{
					"type":        "string",
					"description": "命令作用域，默认为 app",
					"default":     "app",
				},
			},
			"required": []string{"name"},
		},
		func(args map[string]interface{}) (*CallToolResult, error) {
			name, _ := args["name"].(string)
			scope, _ := args["scope"].(string)
			if scope == "" {
				scope = "app"
			}

			cmd := exec.Command("go", "run", ".", "generate:command", "--name", name, "--scope", scope)
			output, err := cmd.CombinedOutput()
			if err != nil {
				return &CallToolResult{
					Content: []ContentBlock{NewTextContent(fmt.Sprintf("生成失败:\n%s", string(output)))},
					IsError: true,
				}, nil
			}
			return &CallToolResult{
				Content: []ContentBlock{NewTextContent(fmt.Sprintf("命令 %s:%s 生成成功:\n%s", scope, name, string(output)))},
			}, nil
		},
	)
}

// registerGenerateTask 注册队列任务生成工具
func registerGenerateTask(s *Server) {
	s.RegisterTool(
		"generate_task",
		"生成异步队列任务。任务会放在 tasks/ 目录下，使用 asynq 作为队列后端。",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "任务名称，例如: send_email, process_order",
				},
			},
			"required": []string{"name"},
		},
		func(args map[string]interface{}) (*CallToolResult, error) {
			name, _ := args["name"].(string)

			cmd := exec.Command("go", "run", ".", "generate:task", "--name", name)
			output, err := cmd.CombinedOutput()
			if err != nil {
				return &CallToolResult{
					Content: []ContentBlock{NewTextContent(fmt.Sprintf("生成失败:\n%s", string(output)))},
					IsError: true,
				}, nil
			}
			return &CallToolResult{
				Content: []ContentBlock{NewTextContent(fmt.Sprintf("任务 %s 生成成功:\n%s", name, string(output)))},
			}, nil
		},
	)
}

// registerInitProject 注册项目初始化工具
func registerInitProject(s *Server) {
	s.RegisterTool(
		"init_project",
		"初始化新的 Lighthouse 项目。会生成完整的项目结构，包括 GraphQL schema、server、models 等。",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "项目名称（Go module 名称），例如: github.com/yourname/myproject",
				},
			},
			"required": []string{"name"},
		},
		func(args map[string]interface{}) (*CallToolResult, error) {
			name, _ := args["name"].(string)

			cmd := exec.Command("go", "run", ".", "generate:init", "--name", name)
			output, err := cmd.CombinedOutput()
			if err != nil {
				return &CallToolResult{
					Content: []ContentBlock{NewTextContent(fmt.Sprintf("初始化失败:\n%s", string(output)))},
					IsError: true,
				}, nil
			}
			return &CallToolResult{
				Content: []ContentBlock{NewTextContent(fmt.Sprintf("项目 %s 初始化成功:\n%s", name, string(output)))},
			}, nil
		},
	)
}

// registerGetDirectiveInfo 注册指令信息查询工具
func registerGetDirectiveInfo(s *Server) {
	s.RegisterTool(
		"get_directive_info",
		"获取 Lighthouse GraphQL 指令的详细用法。支持的指令包括: @loader, @auth, @own, @hidden, @varchar, @text, @longtext, @gorm, @index, @unique, @default 等。",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "指令名称，不带 @ 符号。例如: loader, auth, varchar",
				},
			},
			"required": []string{"name"},
		},
		func(args map[string]interface{}) (*CallToolResult, error) {
			name, _ := args["name"].(string)
			info := getDirectiveInfo(name)
			return &CallToolResult{
				Content: []ContentBlock{NewTextContent(info)},
			}, nil
		},
	)
}

// registerGetConfigInfo 注册配置信息查询工具
func registerGetConfigInfo(s *Server) {
	s.RegisterTool(
		"get_config_info",
		"获取 Lighthouse 配置项的详细说明。支持查询: database, redis, queue, messaging, cors, health 等模块的配置。",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"module": map[string]interface{}{
					"type":        "string",
					"description": "模块名称: database, redis, queue, messaging, cors, health, all",
				},
			},
			"required": []string{"module"},
		},
		func(args map[string]interface{}) (*CallToolResult, error) {
			module, _ := args["module"].(string)
			info := getConfigInfo(module)
			return &CallToolResult{
				Content: []ContentBlock{NewTextContent(info)},
			}, nil
		},
	)
}

// registerListGenerators 注册生成器列表工具
func registerListGenerators(s *Server) {
	s.RegisterTool(
		"list_generators",
		"列出所有可用的代码生成器。",
		map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
		func(args map[string]interface{}) (*CallToolResult, error) {
			generators := `# Lighthouse 代码生成器

## 可用生成器

| 命令 | 说明 |
|------|------|
| generate:init | 初始化新项目 |
| generate:command | 生成 CLI 命令 |
| gen | 生成 GraphQL Schema、Models、DataLoader |
| generate:task | 生成队列任务 |

## 使用方式

### 初始化项目
` + "```bash" + `
lighthouse generate:init --name github.com/yourname/project
` + "```" + `

### 生成 Schema
在项目目录下运行:
` + "```bash" + `
lighthouse gen
# 或
go run . gen
` + "```" + `

### 生成命令
` + "```bash" + `
lighthouse generate:command --name migrate --scope db
` + "```" + `

### 生成任务
` + "```bash" + `
lighthouse generate:task --name send_email
` + "```" + `
`
			return &CallToolResult{
				Content: []ContentBlock{NewTextContent(generators)},
			}, nil
		},
	)
}

// registerSearchDocs 注册文档搜索工具
func registerSearchDocs(s *Server) {
	s.RegisterTool(
		"search_docs",
		"搜索 Lighthouse 框架文档。可以搜索指令、配置、API、最佳实践等内容。",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "搜索关键词",
				},
			},
			"required": []string{"query"},
		},
		func(args map[string]interface{}) (*CallToolResult, error) {
			query, _ := args["query"].(string)
			results := searchDocs(query)
			return &CallToolResult{
				Content: []ContentBlock{NewTextContent(results)},
			}, nil
		},
	)
}

// registerRunLighthouseCommand 注册运行 lighthouse 命令工具
func registerRunLighthouseCommand(s *Server) {
	s.RegisterTool(
		"run_command",
		"运行 lighthouse CLI 命令。可以执行任何 lighthouse 支持的命令。",
		map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"command": map[string]interface{}{
					"type":        "string",
					"description": "要执行的命令，例如: gen, generate:init --name xxx",
				},
			},
			"required": []string{"command"},
		},
		func(args map[string]interface{}) (*CallToolResult, error) {
			command, _ := args["command"].(string)
			parts := strings.Fields(command)

			cmdArgs := append([]string{"run", "."}, parts...)
			cmd := exec.Command("go", cmdArgs...)
			cmd.Dir, _ = os.Getwd()
			output, err := cmd.CombinedOutput()
			if err != nil {
				return &CallToolResult{
					Content: []ContentBlock{NewTextContent(fmt.Sprintf("命令执行失败:\n%s\n错误: %s", string(output), err.Error()))},
					IsError: true,
				}, nil
			}
			return &CallToolResult{
				Content: []ContentBlock{NewTextContent(fmt.Sprintf("命令执行成功:\n%s", string(output)))},
			}, nil
		},
	)
}

// getDirectiveInfo 获取指令信息
func getDirectiveInfo(name string) string {
	directives := map[string]string{
		"loader": `# @loader 指令

用于标记模型启用 DataLoader 批量加载。

## 语法
` + "```graphql" + `
directive @loader(
  keys: [String!]      # 额外的查询键
  morphKey: String     # 多态关联的键字段
  unionTypes: [String!] # 多态关联的类型列表
  extraKeys: [String!] # 额外的索引键
) on OBJECT
` + "```" + `

## 示例
` + "```graphql" + `
type User @loader {
  id: ID!
  name: String!
  posts: [Post!]!
}

# 带额外键
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

## 生成的代码
- 自动生成 DataLoader 批量加载函数
- 支持通过 ID 或指定 keys 批量查询
- 避免 N+1 查询问题
`,
		"auth": `# @auth 指令

用于保护需要认证的字段或查询。

## 语法
` + "```graphql" + `
directive @auth(msg: String) on FIELD_DEFINITION
` + "```" + `

## 示例
` + "```graphql" + `
type Query {
  me: User! @auth
  users: [User!]! @auth(msg: "请先登录")
}

type User {
  id: ID!
  email: String! @auth  # 只有登录用户可以看到邮箱
  name: String!         # 公开字段
}
` + "```" + `

## 工作原理
1. 检查 context 中的用户 ID
2. 如果未登录，返回错误信息
3. 可自定义错误消息
`,
		"own": `# @own 指令

用于所有权校验，确保用户只能访问自己的数据。

## 语法
` + "```graphql" + `
directive @own(field: String!) on FIELD_DEFINITION
` + "```" + `

## 示例
` + "```graphql" + `
type Post {
  id: ID!
  userId: Int!
  content: String!
  secretNotes: String @own(field: "userId")  # 只有帖子作者可以看到
}

type Mutation {
  updatePost(id: ID!, content: String!): Post! @own(field: "userId")
}
` + "```" + `

## 工作原理
1. 获取当前登录用户 ID
2. 与指定字段的值比较
3. 不匹配则返回 nil
`,
		"hidden": `# @hidden 指令

隐藏字段，始终返回 nil。

## 语法
` + "```graphql" + `
directive @hidden on FIELD_DEFINITION
` + "```" + `

## 示例
` + "```graphql" + `
type User {
  id: ID!
  name: String!
  password: String! @hidden  # 永远不返回密码
  internalNote: String @hidden
}
` + "```" + `
`,
		"varchar": `# @varchar 指令

设置字符串字段的 varchar 长度。

## 语法
` + "```graphql" + `
directive @varchar(length: Int!) on FIELD_DEFINITION
` + "```" + `

## 示例
` + "```graphql" + `
type User {
  id: ID!
  name: String! @varchar(length: 100)
  email: String! @varchar(length: 255)
  phone: String @varchar(length: 20)
}
` + "```" + `

## 生成的 GORM tag
` + "```go" + `
Name  string ` + "`gorm:\"type:varchar(100)\"`" + `
Email string ` + "`gorm:\"type:varchar(255)\"`" + `
` + "```" + `
`,
		"text": `# @text 指令

设置字段为 TEXT 类型（最大 65,535 字符）。

## 语法
` + "```graphql" + `
directive @text on FIELD_DEFINITION
` + "```" + `

## 示例
` + "```graphql" + `
type Post {
  id: ID!
  title: String!
  content: String! @text
}
` + "```" + `
`,
		"longtext": `# @longtext 指令

设置字段为 LONGTEXT 类型（最大 4GB）。

## 语法
` + "```graphql" + `
directive @longtext on FIELD_DEFINITION
` + "```" + `

## 示例
` + "```graphql" + `
type Article {
  id: ID!
  title: String!
  body: String! @longtext
}
` + "```" + `
`,
		"gorm": `# @gorm 指令

直接设置 GORM tag。

## 语法
` + "```graphql" + `
directive @gorm(value: String!) on FIELD_DEFINITION
` + "```" + `

## 示例
` + "```graphql" + `
type User {
  id: ID!
  email: String! @gorm(value: "uniqueIndex")
  status: Int! @gorm(value: "default:1")
  metadata: String @gorm(value: "type:json")
}
` + "```" + `
`,
		"index": `# @index 指令

为字段创建数据库索引。

## 语法
` + "```graphql" + `
directive @index(name: String) on FIELD_DEFINITION
` + "```" + `

## 示例
` + "```graphql" + `
type Post {
  id: ID!
  userId: Int! @index
  categoryId: Int! @index(name: "idx_category")
  createdAt: Time! @index
}
` + "```" + `
`,
		"unique": `# @unique 指令

为字段创建唯一索引。

## 语法
` + "```graphql" + `
directive @unique on FIELD_DEFINITION
` + "```" + `

## 示例
` + "```graphql" + `
type User {
  id: ID!
  email: String! @unique
  username: String! @unique
}
` + "```" + `
`,
		"default": `# @default 指令

设置字段的默认值。

## 语法
` + "```graphql" + `
directive @default(value: String!) on FIELD_DEFINITION
` + "```" + `

## 示例
` + "```graphql" + `
type User {
  id: ID!
  status: Int! @default(value: "1")
  role: String! @default(value: "user")
  createdAt: Time! @default(value: "CURRENT_TIMESTAMP")
}
` + "```" + `
`,
	}

	info, ok := directives[name]
	if !ok {
		available := make([]string, 0, len(directives))
		for k := range directives {
			available = append(available, "@"+k)
		}
		return fmt.Sprintf("未找到指令 @%s\n\n可用指令: %s", name, strings.Join(available, ", "))
	}
	return info
}

// getConfigInfo 获取配置信息
func getConfigInfo(module string) string {
	configs := map[string]string{
		"database": `# 数据库配置

Lighthouse 使用 GORM + MySQL，支持主从架构。

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| DB_HOST | localhost | 数据库主机 |
| DB_PORT | 3306 | 数据库端口 |
| DB_USER | root | 数据库用户 |
| DB_PASSWORD | | 数据库密码 |
| DB_NAME | | 数据库名称 |
| DB_MAX_IDLE_CONNS | 10 | 最大空闲连接数 |
| DB_MAX_OPEN_CONNS | 100 | 最大打开连接数 |
| DB_CONN_MAX_LIFETIME | 30 | 连接最大生命周期（分钟） |
| DB_CONN_MAX_IDLE_TIME | 3 | 空闲连接最大存活时间（分钟）**重要** |
| DB_PREPARE_STMT | false | 是否启用 prepared statement 缓存 |

## 主从库配置

| 变量 | 说明 |
|------|------|
| DB_ENABLE_SLAVE | 是否启用从库 |
| DB_MAIN_HOST | 主库地址 |
| DB_SLAVE_HOST | 从库地址（逗号分隔多个） |

## 最佳实践

1. **必须设置 DB_CONN_MAX_IDLE_TIME**：否则空闲连接不会被清理
2. **谨慎使用 DB_PREPARE_STMT=true**：会导致连接累积
3. **监控连接池**：使用 sqlDB.Stats() 监控连接状态
`,
		"redis": `# Redis 配置

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| REDIS_ENABLE | false | 是否启用 Redis |
| REDIS_HOST | localhost | Redis 主机 |
| REDIS_PORT | 6379 | Redis 端口 |
| REDIS_PASSWORD | | Redis 密码 |
| REDIS_DB | 0 | Redis 数据库 |
| REDIS_POOL_SIZE | 10 | 连接池大小 |
| REDIS_MIN_IDLE_CONNS | 5 | 最小空闲连接数 |

## 使用方式

` + "```go" + `
import "github.com/light-speak/lighthouse/redis"

client, err := redis.GetLightRedis()
if err != nil {
    // handle error
}

// 使用 client
client.Set(ctx, "key", "value", 0)
` + "```" + `
`,
		"queue": `# 队列配置

Lighthouse 使用 asynq 作为异步任务队列。

## 环境变量

使用 Redis 配置，队列依赖 Redis 存储任务。

## 使用方式

### 1. 生成任务
` + "```bash" + `
lighthouse generate:task --name send_email
` + "```" + `

### 2. 定义任务
` + "```go" + `
type SendEmailTask struct{}

func (t *SendEmailTask) Name() string {
    return "send_email"
}

func (t *SendEmailTask) Execute(ctx context.Context, payload []byte) error {
    // 处理任务
    return nil
}
` + "```" + `

### 3. 注册任务
` + "```go" + `
queue.RegisterJob("send_email", queue.JobConfig{
    Name:     "send_email",
    Priority: 1,
    Executor: &SendEmailTask{},
})
` + "```" + `

### 4. 投递任务
` + "```go" + `
queue.Enqueue("send_email", payload)
` + "```" + `
`,
		"messaging": `# 消息队列配置

Lighthouse 支持 NATS 消息队列。

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| MESSAGING_DRIVER | nats | 消息驱动（目前只支持 nats） |
| NATS_URL | nats://localhost:4222 | NATS 服务器地址 |

## 使用方式

` + "```go" + `
import "github.com/light-speak/lighthouse/messaging"

broker := messaging.GetBroker()

// 发布消息
broker.Publish("topic", []byte("message"))

// 订阅消息
broker.Subscribe("topic", func(msg []byte) {
    // 处理消息
})
` + "```" + `
`,
		"cors": `# CORS 配置

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| CORS_ALLOW_ORIGINS | * | 允许的来源 |
| CORS_ALLOW_METHODS | GET,POST,OPTIONS | 允许的方法 |
| CORS_ALLOW_HEADERS | Content-Type,Authorization | 允许的请求头 |

## 配置示例

` + "```bash" + `
CORS_ALLOW_ORIGINS=https://example.com,https://app.example.com
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_HEADERS=Content-Type,Authorization,X-Request-ID
` + "```" + `
`,
		"health": `# 健康检查配置

Lighthouse 提供内置的健康检查端点。

## 端点

| 端点 | 说明 |
|------|------|
| /health | 心跳检查 |
| /ready | 就绪检查（检查数据库、内存等） |

## 检查项

1. **数据库连接** - Ping 测试 + 响应时间
2. **内存使用** - 当前内存 vs 阈值（默认 1GB）
3. **连接池使用率** - In-Use / Max（阈值 80%）

## 状态级别

| 状态 | HTTP 状态码 | 含义 |
|------|-------------|------|
| healthy | 200 | 所有检查通过 |
| degraded | 200 | 部分检查失败但可用 |
| unhealthy | 503 | 无法接收流量 |

## 配置

` + "```go" + `
health.SetConfig(&health.Config{
    DBMaxOpenConnsThreshold: 0.8,  // 80% 使用率阈值
    MemoryThresholdMB:       1024, // 1GB 内存阈值
    DBPingTimeout:           3 * time.Second,
})
` + "```" + `
`,
	}

	if module == "all" {
		var result strings.Builder
		result.WriteString("# Lighthouse 配置大全\n\n")
		for name, config := range configs {
			result.WriteString(fmt.Sprintf("---\n\n## %s\n\n", strings.ToUpper(name)))
			result.WriteString(config)
			result.WriteString("\n")
		}
		return result.String()
	}

	info, ok := configs[module]
	if !ok {
		available := make([]string, 0, len(configs))
		for k := range configs {
			available = append(available, k)
		}
		return fmt.Sprintf("未找到模块 %s\n\n可用模块: %s, all", module, strings.Join(available, ", "))
	}
	return info
}

// searchDocs 搜索文档
func searchDocs(query string) string {
	query = strings.ToLower(query)

	// 尝试读取 CLAUDE.md
	claudeMd := ""
	if content, err := os.ReadFile("CLAUDE.md"); err == nil {
		claudeMd = string(content)
	} else {
		// 尝试从 lighthouse 安装目录读取
		if gopath := os.Getenv("GOPATH"); gopath != "" {
			claudePath := filepath.Join(gopath, "pkg/mod/github.com/light-speak/lighthouse@*/CLAUDE.md")
			matches, _ := filepath.Glob(claudePath)
			if len(matches) > 0 {
				if content, err := os.ReadFile(matches[len(matches)-1]); err == nil {
					claudeMd = string(content)
				}
			}
		}
	}

	var results strings.Builder
	results.WriteString(fmt.Sprintf("# 搜索结果: %s\n\n", query))

	// 搜索指令
	directives := []string{"loader", "auth", "own", "hidden", "varchar", "text", "longtext", "gorm", "index", "unique", "default"}
	for _, d := range directives {
		if strings.Contains(d, query) || strings.Contains(query, d) {
			results.WriteString(fmt.Sprintf("## 找到指令: @%s\n\n", d))
			results.WriteString(getDirectiveInfo(d))
			results.WriteString("\n---\n\n")
		}
	}

	// 搜索配置
	configs := []string{"database", "redis", "queue", "messaging", "cors", "health"}
	for _, c := range configs {
		if strings.Contains(c, query) || strings.Contains(query, c) {
			results.WriteString(fmt.Sprintf("## 找到配置: %s\n\n", c))
			results.WriteString(getConfigInfo(c))
			results.WriteString("\n---\n\n")
		}
	}

	// 如果有 CLAUDE.md，搜索其中的内容
	if claudeMd != "" && strings.Contains(strings.ToLower(claudeMd), query) {
		results.WriteString("## 在项目文档中找到相关内容\n\n")
		lines := strings.Split(claudeMd, "\n")
		for i, line := range lines {
			if strings.Contains(strings.ToLower(line), query) {
				start := i - 2
				if start < 0 {
					start = 0
				}
				end := i + 5
				if end > len(lines) {
					end = len(lines)
				}
				results.WriteString("```\n")
				results.WriteString(strings.Join(lines[start:end], "\n"))
				results.WriteString("\n```\n\n")
			}
		}
	}

	if results.Len() == len(fmt.Sprintf("# 搜索结果: %s\n\n", query)) {
		results.WriteString("未找到相关内容。\n\n")
		results.WriteString("建议:\n")
		results.WriteString("- 尝试搜索指令名称: loader, auth, varchar 等\n")
		results.WriteString("- 尝试搜索配置模块: database, redis, queue 等\n")
		results.WriteString("- 使用 get_directive_info 或 get_config_info 工具获取详细信息\n")
	}

	return results.String()
}
