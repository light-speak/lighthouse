# 数据库迁移

Lighthouse 使用 [Atlas](https://atlasgo.io/) 进行数据库迁移管理。

## Atlas 配置

项目初始化时会生成 `atlas.hcl` 配置文件：

```hcl
data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./loader",
  ]
}

env "dev" {
  src = data.external_schema.gorm.url
  dev = "mysql://root:@127.0.0.1:3306/test"
  url = "mysql://root:@127.0.0.1:3306/myapp"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "production" {
  src = data.external_schema.gorm.url
  url = env("DATABASE_URL")
  migration {
    dir = "file://migrations"
  }
}
```

## GORM Schema 加载器

`loader/main.go` 用于将 GORM 模型转换为 Atlas schema：

```go
package main

import (
    "myapp/models"
    "ariga.io/atlas-provider-gorm/gormschema"
    "gorm.io/gorm"
)

var migrateModels = []interface{}{
    &models.User{},
    &models.Wallet{},
    &models.Post{},
    // 添加所有需要迁移的模型...
}

func main() {
    option := gormschema.WithConfig(&gorm.Config{
        DisableForeignKeyConstraintWhenMigrating: true,
        IgnoreRelationshipsWhenMigrating:         true,
    })
    stmts, _ := gormschema.New("mysql", option).Load(migrateModels...)
    io.WriteString(os.Stdout, stmts)
}
```

::: tip
每次添加新模型后，记得更新 `migrateModels` 列表。
:::

## 迁移工作流

### 1. 修改 Schema

编辑 `schema/*.graphql` 文件：

```graphql
type User @loader(keys: ["id"]) {
  id: ID!
  name: String! @varchar(length: 100)
  email: String @unique  # 新增字段
}
```

### 2. 生成 Go 代码

```bash
lighthouse generate:schema
```

### 3. 更新 Loader

编辑 `loader/main.go`，添加新模型（如果有新模型）。

### 4. 生成迁移文件

```bash
atlas migrate diff --env dev
```

这会在 `migrations/` 目录生成迁移 SQL 文件：

```
migrations/
├── 20240101000000.sql
├── 20240102000000.sql
└── atlas.sum
```

### 5. 检查迁移 SQL

查看生成的 SQL 确保正确：

```sql
-- migrations/20240102000000.sql
ALTER TABLE `users` ADD COLUMN `email` varchar(255) UNIQUE;
```

### 6. 应用迁移

```bash
go run . migration:apply --env=dev
```

## 常用命令

```bash
# 生成迁移
atlas migrate diff --env dev

# 应用迁移
go run . migration:apply --env=dev

# 查看迁移状态
atlas migrate status --env dev

# 回滚迁移（谨慎使用）
atlas migrate down --env dev
```

## 生产环境迁移

```bash
# 设置环境变量
export DATABASE_URL="mysql://user:password@host:3306/dbname"

# 应用迁移
go run . migration:apply --env=production
```

## 注意事项

::: warning 外键约束
默认禁用外键约束迁移，避免复杂的依赖关系：
```go
gormschema.WithConfig(&gorm.Config{
    DisableForeignKeyConstraintWhenMigrating: true,
})
```
:::

::: danger 生产环境
- 在生产环境应用迁移前，先在测试环境验证
- 备份数据库
- 在低峰期执行迁移
:::
