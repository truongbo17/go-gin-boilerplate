# Architecture and extension points

The application builds dependencies once at startup. HTTP requests reuse services and connection pools. Transport validates input, core applies rules, repositories persist data, and infrastructure packages provide reusable clients.

## Package map

| Package | Responsibility |
| --- | --- |
| `cmd`, `cmd/cli` | Process commands, configuration loading, resource lifetime, database CLI commands |
| `config` | Environment variables and startup validation |
| `internal/routes` | Gin engine, shared middleware, public routes, GraphQL mount, OpenAPI |
| `internal/app/v1/<feature>` | HTTP feature module, REST routes, controllers, requests, responses |
| `internal/app/graphql/<feature>` | GraphQL schema, generated code, resolvers, query limits |
| `internal/app/core/graphql`, `internal/repository/graphql` | Enum catalog, entity option policy, allowlisted metadata queries |
| `internal/app/core/<feature>` | Use cases, feature data, errors, interfaces consumed by services |
| `internal/repository/<feature>` | Queries and transactions for that feature |
| `internal/infra` | Reusable database, cache, messaging, storage, logging, tracing, worker adapters |
| `internal/middlewares` | Shared HTTP middleware |
| `internal/app/worker/<feature>` | Task transport adapters for a feature |
| `internal/app/worker/register` | Handler and schedule assembly for one worker process |
| `internal/app/worker/app.go` | Enabled worker feature list and handler registry construction |
| `internal/app/core/worker` | Stable task names, parameter types, and dispatch contract |

`internal/request`, `internal/response`, `internal/i18n`, and `internal/page` contain shared HTTP, localization, and pagination support. Put only small, domain-independent parsing helpers in `internal/utils`.

## Construction and request flow

```text
cmd/server.go
  -> cmd/runtime.go opens MySQL, optional Redis, cache, logger, tracer
  -> internal/routes.New builds Gin and shared middleware
       -> constructs enabled features and mounts each on /api

HTTP request
  -> shared middleware -> route rate limit -> JWT/permission -> validator
  -> controller -> core service -> feature repository -> database
```

`cmd/runtime.go` owns process resources. `cmd/server.go` passes opened resources once to `internal/routes.New`, which constructs enabled features and registers their routes in order. GraphQL receives auth middleware from the already built auth module. Controllers, task handlers, and core services keep their specific dependencies. Adding an auth endpoint normally uses the existing auth module fields. A new feature gets its own module instead of enlarging the auth module.

Services declare small interfaces for the operations they consume. The auth service owns `AuthUserRepository` and `TokenBlacklist`; `internal/repository/auth` and `internal/infra/cache` satisfy them. Core services use `services.AuthSettings`; the HTTP module maps values from `config.Auth` at startup. Keep GORM queries and transactions in repositories. Map core errors to HTTP status and localized messages in the feature response package.

The bundled auth models reuse `internal/model`, including GORM tags and `gorm.DeletedAt`. This is an intentional persistence coupling for the MySQL-backed auth base. PostgreSQL is not a drop-in replacement for auth. If a feature needs independent domain and storage representations, add mapping in its repository then.

## Transport ownership

REST route definitions for auth live in `internal/app/v1/auth/routes.go`, beside the module and controllers. Each route shows its rate limit, authentication, permission check, validator, and handler in execution order. `internal/routes` owns shared Gin setup, public endpoints, and Swagger/OpenAPI; documentation routes are registered only outside `release`. The GraphQL HTTP module under `internal/app/v1/graphql` owns its routes; resolvers and generated schema remain under `internal/app/graphql/auth`.

HTTP request types and validation belong in the feature's `requests` package. Reusable Gin binding and path parsing belong in `internal/request`. Core input types describe use cases and contain no HTTP response behavior. Repositories return feature data; `internal/response` builds HTTP pagination links from `internal/page` results.

## Process and resource lifetime

`cmd/runtime.go` loads fresh configuration for each command and closes resources in reverse order on shutdown or startup failure. `server` opens MySQL, cache, and optional tracing; it also creates an Asynq producer when mail is enabled. `worker` opens Redis and optional tracing, plus MySQL when mail is enabled. `migrate` and `create_user` open MySQL; `version` opens nothing. CLI definitions live in `cmd/cli` and use the runtime's database runner rather than opening another pool.

`internal/app/core/auth/services.MailService` owns recipient lookup, mail content, and reset-link generation. Auth services submit typed jobs through `internal/app/core/worker.Dispatcher`; the worker auth handler decodes each task and calls `MailService`. `internal/app/worker/app.go` lists enabled worker features; each handler lists its own task types in `Handlers()`, and `register.New` combines them while rejecting duplicate types. `cmd/worker.go` calls the composition function once and owns process startup and shutdown. There is no mutable registry in core. Optional PostgreSQL, MongoDB, ClickHouse, Kafka, and S3 clients are available in `internal/infra`; the default server does not connect to them. A feature that needs one should acquire it during startup, reuse it, and arrange shutdown with the owning command.

## Extend the application

### REST feature

For an existing feature, define transport-independent input/output and rules in `internal/app/core/<feature>`, persistence in `internal/repository/<feature>`, then HTTP requests, controllers, responses, and routes in `internal/app/v1/<feature>`. Each route keeps its own rate limit, authentication, permission, validator, and handler in execution order. Update `internal/routes/openapi.yaml` when the HTTP contract changes. For a new feature, follow auth's `Dependencies`, `NewModule`, and `RegisterRoutes`; construct and mount it in `internal/routes/main.go`. Give each controller only the services it uses. Introduce a shared package only when more than one feature needs the behavior.

### GraphQL

Edit `internal/app/graphql/auth/schema.graphqls`, implement its resolver beside the schema, and regenerate with `go run github.com/99designs/gqlgen generate --config gqlgen.yml`. Resolvers adapt input and output; `internal/app/core/graphql` owns the enum catalog and entity access policy, while `internal/repository/graphql` owns allowlisted database queries. `internal/app/v1/graphql` registers the public enum endpoint, authenticated GraphQL endpoint, and authenticated REST entity-options endpoint. Entity keys `users`, `roles`, and `permissions` require their matching RBAC read permissions. GraphQL uses a 64 KiB body cap, 2-second request timeout, parser token/depth/alias/fragment/complexity limits, and disables introspection in release mode. Add only entity keys and enum groups owned by this application; the ERP export is a design reference, not a source of domain data.

### Worker task and schedule

Define stable task names and minimal payloads in `internal/app/core/worker`. Core services dispatch tasks by type. `internal/app/worker/<feature>` decodes and validates tasks, then calls the feature service. Follow `internal/app/worker/auth`: create `Dependencies`, `NewHandler`, and `Handlers() map[string]asynq.HandlerFunc`, then add that handler to `internal/app/worker/app.go` when enabled. `internal/app/worker/register` combines handlers without knowing feature names or business logic. Put schedules in `register/schedule.go` and let `cmd/worker.go` own startup and shutdown. Use `asynq.SkipRetry` for malformed or permanent failures; return transient errors for retry. Make handlers safe for repeated delivery. The auth mail example queues no password or reset token; `MailService` creates the token when handling the job. The example schedule prints every 15 seconds while the worker runs.

Queued tasks survive a worker restart while Redis retains data. Transient mail failures retry five times, then archive. Registration can commit before enqueue fails because this example has no transactional outbox. SMTP acceptance can be ambiguous after a connection failure, so email may be duplicated. Add an outbox and durable Redis setup when delivery guarantees require recovery across these failures.

### CLI and migrations

Add Cobra commands in `cmd/cli` and register them in `cmd/root.go`. Use `cmd/runtime.go` for resources; database-only commands use `cli.WithDatabase` to avoid importing the parent `cmd` package. Validate input before opening a connection. Add migrations with unique IDs under `internal/migrations` and register their order in `internal/migrations/main.go`. The bundled auth migrations target MySQL.

## Optional infrastructure

Open only the clients needed by an enabled feature, reuse them, and close them with the owning command. Constructors return clients and errors rather than writing package globals. Feature-specific queries and topic or object naming rules belong outside generic `internal/infra` clients.

| Backend | Constructor | Lifecycle and scope |
| --- | --- | --- |
| PostgreSQL | `database.OpenPostgres(ctx, opts)` | Close the GORM SQL pool; auth migrations require a separate PostgreSQL implementation. |
| MongoDB | `mongodb.Open(ctx, opts)` | Disconnect on shutdown; put document queries in a feature repository. |
| ClickHouse | `clickhouse.Open(ctx, opts)` | Close the native connection. |
| Kafka | `kafka.NewWriter` / `kafka.NewReader` | Close writer/reader; commit consumed messages after successful handling. |
| S3-compatible storage | `objectstore.OpenS3(ctx, opts)` | Uses the AWS credential chain; check bucket access when the feature requires it. |
| SMTP | `mail.Sender{Config: cfg.Mail}` | Opens a connection per send; TLS modes verify certificates. |

Kafka consumers and Asynq workers can deliver more than once. Keep their handlers idempotent.

## Local performance reference

The 2026-10-08 local baseline used native Go 1.26.8 binaries on a 10-core Apple Silicon Mac, MySQL 8.0.46 in Docker, loopback HTTP/1.1, local cache, no tracing, and a same-host load generator. Each public-endpoint run had a 2-second warmup and 10 seconds of measured traffic. These numbers are a comparison point, not a production capacity estimate.

| Endpoint | Connections | Requests/s | p99 |
| --- | ---: | ---: | ---: |
| `/ping` | 16 | 104,081 | 0.516 ms |
| `/ping` | 64 | 144,281 | 1.631 ms |
| `/ping` | 128 | 151,630 | 3.454 ms |

The load generator and server shared CPU, and this run did not isolate a database or CPU bottleneck. Reproduce with `scripts/bench_http.go` on representative infrastructure before making capacity decisions.

## Design rationale

The feature-module boundary keeps route signatures short. `cmd/runtime.go` owns resource lifetimes, while `internal/routes` owns shared Gin setup and feature registration. Required resources are checked when the router starts. Auth models currently retain GORM metadata by design; use separate domain and storage types only when a feature needs that boundary.

## References

- [Go: Package names](https://go.dev/blog/package-names)
- [Go: Dependency injection with Wire](https://go.dev/blog/wire) — describes explicit constructor wiring, the pattern used here without code generation.
- [Kratos project layout](https://github.com/go-kratos/kratos-layout) — a reference for transport, business, and data separation; this repository keeps fewer packages.
