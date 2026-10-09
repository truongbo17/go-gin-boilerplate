# Coding conventions

These conventions describe the code in this repository. Use the Go version declared in `go.mod`, run `gofmt`, and follow [Effective Go](https://go.dev/doc/effective_go) and [Go code review comments](https://go.dev/wiki/CodeReviewComments).

## Packages and dependencies

- Keep package names short, lower case, and meaningful at the import site. Name identifiers `ID`, `URL`, and `HTTP` in new code.
- Keep use-case decisions in `internal/app/core/<feature>`. Core services do not import Gin, HTTP response packages, application config, or infrastructure connection packages. The current auth models intentionally retain GORM metadata; see [Architecture](architecture.md).
- Put feature persistence in `internal/repository/<feature>`, shared clients in `internal/infra`, route middleware in `internal/middlewares`, task adapters and their task-type mapping in `internal/app/worker/<feature>`, and per-process registration in `internal/app/worker/register`.
- Put feature REST routes, controllers, requests, responses, and module wiring under `internal/app/v1/<feature>`. Keep the top-level router in `internal/routes`; put GraphQL HTTP routes in `internal/app/v1/graphql`, schemas and resolvers in `internal/app/graphql/<feature>`, and feature-neutral metadata rules in `internal/app/core/graphql`.
- Do not add a package-wide DB, Redis client, service locator, or mutable job registry. `cmd/runtime.go` owns connections; `internal/routes/main.go` constructs enabled HTTP features and calls their `RegisterRoutes` methods directly.
- Define an interface near the service that consumes it when it clarifies a boundary. Prefer a concrete type for a single implementation with no such need.

## Functions and state

- Pass `context.Context` as the first parameter of request-scoped service and repository methods. Use the same context for database, cache, queue, and outbound HTTP work.
- Return errors with operation context using `%w`. Handle startup, query, migration, and shutdown errors. Keep HTTP status and public error messages out of repositories and services.
- Give constructors only the dependencies they use. Group related constructor inputs in a feature-owned `Dependencies` struct when a call would otherwise need many parameters. Pass opened resources once to `internal/routes.New`, validate required resources, and avoid passing `*runtime` into services, controllers, or repositories.
- Keep process resource ownership in `cmd/runtime.go`. A constructor returning a pool or client needs a matching close path there. Do not open connections per request.
- Validate external input at entry points. Transport validators handle HTTP input; CLI commands validate flags before opening resources. Core services still enforce rules shared by transports.

## HTTP routes

Put route middleware on separate lines so policy is visible during review:

```go
rbac.GET("/roles",
    module.rateLimit("role:list", 300),
    jwt,
    middlewares.CheckPermission("role:index", module.permissionService),
    requests.ListRoleValidator(),
    module.roleController.ListRoles,
)
```

Use distinct rate-limit policy names for endpoints with independent counters. Authentication runs before permission checks; validation runs before controllers. The top-level router owns security headers, body limits, request IDs, logging, CORS, and trusted-proxy settings. Update `internal/routes/openapi.yaml` when a REST contract changes.

## Data and external services

- Keep GORM `Where`, `Joins`, and transactions in repositories. Use parameterized conditions. Reuse `internal/page` for list results; build URLs only in the HTTP response layer.
- Construct optional clients only for features that need them. Place feature-specific MongoDB queries, Kafka topic rules, and S3 object naming outside the generic clients in `internal/infra`.
- Use Redis-backed cache and rate limits for multiple server instances. Local cache suits development or a single instance.

## Changes and verification

- Run `gofmt`, `go vet ./...`, `go build ./...`, and `git diff --check` before review. For database, queue, or migration changes, run the affected command against disposable local services and inspect the result.
- Update README, `docs/architecture.md`, `docs/go-style.md`, OpenAPI, or the GraphQL schema when a command, route, package boundary, or public contract changes.
