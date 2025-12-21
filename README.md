<p align="center">
  <img src="docs/public/logo.svg" width="80" height="80" alt="Lighthouse">
</p>

<h1 align="center">Lighthouse</h1>

<p align="center">
  <strong>Build GraphQL APIs in Go. Fast.</strong>
</p>

<p align="center">
  <a href="https://github.com/light-speak/lighthouse/releases"><img src="https://img.shields.io/github/v/release/light-speak/lighthouse?style=flat-square&color=blue" alt="Release"></a>
  <a href="https://pkg.go.dev/github.com/light-speak/lighthouse"><img src="https://pkg.go.dev/badge/github.com/light-speak/lighthouse.svg" alt="Go Reference"></a>
  <a href="https://goreportcard.com/report/github.com/light-speak/lighthouse"><img src="https://goreportcard.com/badge/github.com/light-speak/lighthouse" alt="Go Report Card"></a>
  <a href="https://github.com/light-speak/lighthouse/blob/main/LICENSE"><img src="https://img.shields.io/github/license/light-speak/lighthouse?style=flat-square" alt="License"></a>
</p>

<p align="center">
  <a href="https://light-speak.github.io/lighthouse/">Documentation</a> •
  <a href="https://light-speak.github.io/lighthouse/guide/getting-started">Quick Start</a> •
  <a href="./README_CN.md">中文</a>
</p>

---

## What is Lighthouse?

Lighthouse is a **batteries-included** GraphQL framework for Go. Define your schema, run one command, and get a production-ready API with DataLoaders, authentication, and database migrations.

```graphql
type User @loader(keys: ["id"]) {
  id: ID!
  name: String! @varchar(length: 100)
  posts: [Post!]! @goField(forceResolver: true)
}

extend type Query {
  me: User! @auth
}
```

```bash
lighthouse generate:schema  # That's it. Models, resolvers, dataloaders generated.
```

## Features

| Feature | Description |
|---------|-------------|
| **Schema-First** | Define GraphQL schema, generate Go code |
| **DataLoader** | Auto-generated, N+1 problem solved |
| **Auth Directives** | `@auth`, `@own` built-in |
| **Database** | GORM + MySQL, connection pooling, master-slave |
| **Migrations** | Atlas-powered schema migrations |
| **Queue** | Redis-based async jobs (asynq) |
| **Messaging** | NATS pub/sub for real-time |
| **Storage** | S3/MinIO/COS unified interface |
| **Metrics** | Prometheus + health checks |

## 5-Minute Start

```bash
# Install
go install github.com/light-speak/lighthouse@latest

# Create project
lighthouse generate:init --module github.com/you/myapp --models user,post
cd myapp

# Configure .env, then run
go run . app:start
```

Open http://localhost:8080 → GraphQL Playground ready.

## Project Structure

```
myapp/
├── schema/          # GraphQL definitions
├── models/          # Generated Go structs
├── resolver/        # Your business logic
├── graph/           # gqlgen generated (don't touch)
├── commands/        # CLI commands
├── server/          # HTTP server setup
└── migrations/      # Database migrations
```

## Tech Stack

- **[gqlgen](https://gqlgen.com/)** - GraphQL engine
- **[GORM](https://gorm.io/)** - ORM
- **[Atlas](https://atlasgo.io/)** - Migrations
- **[asynq](https://github.com/hibiken/asynq)** - Job queue
- **[NATS](https://nats.io/)** - Messaging
- **[zerolog](https://github.com/rs/zerolog)** - Logging

## Documentation

📚 **[Full Documentation](https://light-speak.github.io/lighthouse/)**

- [Getting Started](https://light-speak.github.io/lighthouse/guide/getting-started)
- [Schema Basics](https://light-speak.github.io/lighthouse/schema/basics)
- [DataLoader](https://light-speak.github.io/lighthouse/schema/dataloader)
- [Authentication](https://light-speak.github.io/lighthouse/features/auth)
- [Database](https://light-speak.github.io/lighthouse/features/database)

## Contributing

PRs welcome. Open an issue first for major changes.

## License

MIT

---

<p align="center">
  <sub>Built with ☕ by <a href="https://github.com/light-speak">Light Speak</a></sub>
</p>
