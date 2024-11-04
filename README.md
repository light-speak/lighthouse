# **🚢 Lighthouse GraphQL Framework**

[English](https://github.com/light-speak/lighthouse/blob/main/README.md) | [中文](https://github.com/light-speak/lighthouse/blob/main/README_zh.md)

[![CI](https://github.com/light-speak/lighthouse/actions/workflows/main.yml/badge.svg)](https://github.com/light-speak/lighthouse/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/light-speak/lighthouse/branch/main/graph/badge.svg)](https://codecov.io/gh/light-speak/lighthouse)
[![Go Report Card](https://goreportcard.com/badge/github.com/light-speak/lighthouse)](https://goreportcard.com/report/github.com/light-speak/lighthouse)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

`Lighthouse` is a feature-rich self-developed GraphQL framework designed to simplify GraphQL service development based on microservice architecture. The framework integrates a logging system (using `zeroLog`), supports Elasticsearch, file logging mode and Redis caching, while featuring a flexible directive system and powerful custom configuration capabilities. The framework currently has built-in support for `gorm` and will expand to more ORM options in the future.

## Features

- **Microservice Architecture Support**: Adopts independent microservice mode, doesn't support GraphQL Federation but manages microservices through a custom service registry.
- **Custom Directives**: Supports rich custom directives for dynamic queries, filtering, relationships, and more.
- **Extensibility**: Supports various GraphQL file structures through configuration files to flexibly meet different project needs.
- **ORM Integration**: Currently supports `gorm`, with plans to add support for more ORM libraries.
- **Logging and Cache Integration**: Integrates `zeroLog` logging system, supports Elasticsearch, file logging, and Redis caching.

## Quick Start

### Installation

1. **Install using `go install`**

   ```bash
   go install github.com/light-speak/lighthouse@latest
   ```

2. **Create a new project**

   ```bash
   lighthouse generate:init
   ```

### Configuration File (lighthouse.yml)

`lighthouse.yml` is the core configuration file for `lighthouse`, used to specify GraphQL Schema paths, file extensions, ORM settings, etc. Here's an example configuration:

```yaml
# lighthouse.yml

schema:
  ext:
    - graphql       # Supported file extensions
    - graphqls
  path:
    - schema        # Path to GraphQL Schema files
  model:
    orm: gorm       # ORM configuration, currently supports gorm
```

- `schema.ext`: Specifies Schema file extensions, can be `.graphql` or `.graphqls`.
- `schema.path`: Defines the path to Schema files, framework will automatically load all files in this path.
- `model.orm`: Currently supports `gorm` as the ORM library.

### Directory Structure

The `example` project structure is as follows:

```plaintext
.
├── cmd                     # CLI related code
│   ├── cmd.go              # Main command entry
│   ├── migrate
│   │   └── migrate.go      # Database migration logic
│   └── start
│       └── start.go        # Service start entry
├── models                  # Data model definitions
│   ├── enum.go             # Enum type definitions
│   ├── input.go            # Input type definitions
│   ├── interface.go        # Interface definitions
│   ├── model.go            # Model structures
│   └── response.go         # Response data structures
├── repo                    # Database operation encapsulation
│   └── repo.go
├── resolver                # GraphQL resolvers
│   ├── mutation.go         # Mutation resolver
│   ├── query.go            # Query resolver
│   └── resolver.go         # Resolver main entry
├── schema                  # GraphQL Schema files
│   └── user.graphql        # Example Schema file
└── service                 # Service logic
    └── service.go
```

### Next Steps

Add your custom schema files in the `schema` directory, then run the following command to generate corresponding code:

```bash
lighthouse generate:schema
```

### Using Directives

In `lighthouse`, you can use the following directives in your GraphQL Schema:

- **@skip / @include**: Conditional query directives for dynamically controlling field inclusion in responses.
- **@enum**: For defining enum type fields, currently only supports `int8` type.
- **@paginate / @find / @first**: For pagination, finding, and getting first result queries.
- **@in / @eq / @neq / @gt / @gte / @lt / @lte / @like / @notIn**: Parameter filtering directives supporting various comparison operators.
- **@belongsTo / @hasMany / @hasOne / @morphTo / @morphToMany**: Relationship mapping directives for defining model relationships.
- **@index / @unique**: Create index or add unique constraint for fields.
- **@defaultString / @defaultInt**: Set default values for fields.
- **@tag**: For marking additional field attributes.
- **@model**: Mark type as database model.
- **@softDeleteModel**: Mark type as database model with soft delete support.
- **@order**: For sorting query results.
- **@cache**: For caching query results to improve response speed.

### Example Code

Here's an example query using the `@paginate` directive for user data pagination:

```graphql
type Query {
  users: [User] @paginate(scopes: ["active"])
}

type User @model(name: "UserModel") {
  id: ID!
  name: String!
  age: Int
  posts: [Post] @hasMany(relation: "Post", foreignKey: "user_id")
}
```

## Extension and Customization

`lighthouse` provides flexible extension interfaces, you can:

- **Add Custom Directives**: Write your own directives to extend framework functionality.
- **Support Other ORMs**: Add support for other ORM libraries by referencing the `gorm` integration approach.

## Development Plan

| 🚀 Feature Category | ✨ Feature Description | 📅 Status |
|-------------------|---------------------|-----------|
| 🛠️ Custom Directives | Add support for custom directives | ✅ Completed |
| 📊 Query Directives | Add @find and @first annotations for query support | ✅ Completed |
| 🔍 Query & Filtering | Add date range filtering directives | ✅ Completed |
|                    | Add string matching directives | ✅ Completed |
|                    | Add dynamic sorting directives | ✅ Completed |
| 📊 Pagination | Add @paginate annotation for pagination support | ✅ Completed |
| 📜 Conditional Query | Add @skip and @include conditional query directives | 🚧 In Progress |
| 📚 Relationship Mapping | Add @morphTo, @morphToMany, @hasOne, @manyToMany | 🚧 In Progress |
| 🔧 Microservice Management | Add microservice registry | ⏳ Planned |
| 💾 Cache Integration | Integrate Redis as cache support | ✅ Completed |
| 📝 Logging System | Integrate zeroLog system, support Elasticsearch and file logging | ✅ Completed |
| 🔄 Cache Directive | Add @cache directive to support query result caching | 🚧 In Progress |
| 🔀 Sorting Directive | Add @order directive to support query result sorting | 🚧 In Progress |
| 🗄️ ORM Support | Extend support for other ORMs like `ent`, `sqlc` | ⏳ Planned |
| 📑 Doc Generation | Auto-generate GraphQL Schema documentation | ⏳ Planned |
| 📦 Plugin Support | Provide plugin system for community contributions | ⏳ Planned |
| 🌐 Frontend Tools | Develop Apollo Studio-like frontend for query testing | ⏳ Planned |
| 📊 Performance Tracking | Support performance tracking for fields and services | ⏳ Planned |

