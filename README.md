<p align="center">
  <h1 align="center">Lighthouse</h1>
  <p align="center">A powerful GraphQL framework for Go with batteries included</p>
</p>

<p align="center">
  <a href="https://github.com/light-speak/lighthouse/releases"><img src="https://img.shields.io/github/v/release/light-speak/lighthouse?style=flat-square" alt="Release"></a>
  <a href="https://pkg.go.dev/github.com/light-speak/lighthouse"><img src="https://pkg.go.dev/badge/github.com/light-speak/lighthouse.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/light-speak/lighthouse"><img src="https://goreportcard.com/badge/github.com/light-speak/lighthouse" alt="Go Report Card"></a>
  <a href="https://github.com/light-speak/lighthouse/blob/main/LICENSE"><img src="https://img.shields.io/github/license/light-speak/lighthouse?style=flat-square" alt="License"></a>
  <a href="https://github.com/light-speak/lighthouse/stargazers"><img src="https://img.shields.io/github/stars/light-speak/lighthouse?style=flat-square" alt="Stars"></a>
</p>

<p align="center">
  <a href="./README_CN.md">中文文档</a>
</p>

---

## Overview

**Lighthouse** is a full-featured GraphQL framework for Go, inspired by [Laravel Lighthouse](https://lighthouse-php.com/). It provides an elegant and efficient way to build GraphQL APIs with built-in support for database management, caching, message queuing, and more.

## Features

- **GraphQL First** - Built on top of [gqlgen](https://gqlgen.com/), the most popular GraphQL library for Go
- **Database Management** - GORM-based MySQL support with connection pooling, master-slave replication
- **DataLoader** - Built-in DataLoader pattern for efficient N+1 query prevention
- **Redis Integration** - Connection pooling, caching utilities, and pub/sub support
- **Async Queue** - Background job processing powered by [asynq](https://github.com/hibiken/asynq)
- **Messaging** - NATS-based real-time messaging and event broadcasting
- **Authentication** - JWT-based authentication with GraphQL directives
- **Code Generation** - CLI tool for scaffolding models, resolvers, and DataLoaders
- **Graceful Shutdown** - Proper resource cleanup on application termination

## Requirements

- Go 1.24+
- MySQL 5.7+ or 8.0+
- Redis 6.0+ (optional)
- NATS 2.0+ (optional)

## Installation

```bash
go get github.com/light-speak/lighthouse
```

### Install CLI Tool

```bash
go install github.com/light-speak/lighthouse@latest
```

## Quick Start

### 1. Initialize a New Project

```bash
lighthouse generate:init --module=github.com/yourname/myproject --models=user,post
cd myproject
```

### 3. Configure Environment

Edit `.env` file with your database credentials:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=myproject
```

### 4. Define GraphQL Schema

Edit `graph/schema.graphqls`:

```graphql
type User {
  id: ID!
  name: String!
  email: String!
  createdAt: DateTime!
}

type Query {
  users: [User!]!
  user(id: ID!): User
}

type Mutation {
  createUser(name: String!, email: String!): User!
}
```

### 5. Generate Code

```bash
lighthouse generate:schema
```

### 6. Run the Server

```bash
go run . app:start
```

Your GraphQL server is now running at `http://localhost:8080/graphql`

## Configuration

### Database Connection Pool

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_MAX_IDLE_CONNS` | 10 | Maximum idle connections |
| `DB_MAX_OPEN_CONNS` | 100 | Maximum open connections |
| `DB_CONN_MAX_LIFETIME` | 30 | Connection max lifetime (minutes) |
| `DB_CONN_MAX_IDLE_TIME` | 3 | Idle connection max lifetime (minutes) |
| `DB_PREPARE_STMT` | false | Enable prepared statement cache |

### Redis Connection Pool

| Variable | Default | Description |
|----------|---------|-------------|
| `REDIS_ENABLE` | false | Enable Redis |
| `REDIS_HOST` | localhost | Redis host |
| `REDIS_PORT` | 6379 | Redis port |
| `REDIS_PASSWORD` | | Redis password |
| `REDIS_DB` | 0 | Redis database |
| `REDIS_POOL_SIZE` | 10 | Connection pool size |
| `REDIS_MIN_IDLE_CONNS` | 5 | Minimum idle connections |

See [.env.example](./lightcmd/initization/tpl/env.tpl) for full configuration options.

## Project Structure

```
lighthouse/
├── databases/       # Database connection management (GORM + MySQL)
├── redis/           # Redis client management
├── queue/           # Async task queue (asynq)
├── messaging/       # Messaging system (NATS)
├── routers/         # GraphQL routers and middleware
│   ├── auth/        # Authentication directives
│   ├── dataloader/  # DataLoader batch queries
│   └── health/      # Health check endpoints
├── lightcmd/        # CLI and code generation
│   ├── generate/    # Code generators
│   └── initization/ # Project initialization templates
├── logs/            # Logging module
├── storages/        # Storage adapters (S3/COS)
├── templates/       # Template engine
├── utils/           # Utility functions
└── lighterr/        # Error handling
```

## CLI Commands

```bash
# Show all available commands
lighthouse help

# Initialize new project
lighthouse generate:init --module=<module-name> --models=<model1,model2>

# Generate schema (models, resolvers, dataloaders)
lighthouse generate:schema

# Generate a new command
lighthouse generate:command --name=<command-name>

# Generate a new task
lighthouse generate:task --name=<task-name>

# Initialize queue service
lighthouse queue:init
```

## DataLoader Pattern

Lighthouse automatically generates DataLoaders to prevent N+1 query problems:

```go
// Auto-generated DataLoader usage
user, err := GetUserIdLoader(ctx).Load(ctx, userID)

// Batch loading
users, err := GetUserIdLoader(ctx).LoadAll(ctx, userIDs)
```

## Authentication

Use GraphQL directives for authentication:

```graphql
directive @auth on FIELD_DEFINITION

type Query {
  me: User! @auth
  publicData: String!
}
```

## Health Checks

Built-in health check endpoints for Kubernetes:

- **Liveness**: `GET /health`
- **Readiness**: `GET /ready`

## Versioning

This project follows [Semantic Versioning](https://semver.org/).

| Version | Status | Go Version |
|---------|--------|------------|
| v1.1.x | Current | 1.24+ |
| v1.0.x | Maintenance | 1.21+ |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [gqlgen](https://gqlgen.com/) - GraphQL server library for Go
- [GORM](https://gorm.io/) - ORM library for Go
- [Laravel Lighthouse](https://lighthouse-php.com/) - Inspiration for this project
- [asynq](https://github.com/hibiken/asynq) - Async task processing
- [NATS](https://nats.io/) - Messaging system

## Support

- [Documentation](https://github.com/light-speak/lighthouse/wiki)
- [Issues](https://github.com/light-speak/lighthouse/issues)
- [Discussions](https://github.com/light-speak/lighthouse/discussions)

---

<p align="center">Made with ❤️ by <a href="https://github.com/light-speak">Light Speak</a></p>
