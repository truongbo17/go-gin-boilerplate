# Go Gin Boilerplate

A reusable Gin API starter. It includes JWT authentication, role-based access control, MySQL migrations, optional Redis cache and workers, rate limiting, structured logs, request IDs, and optional OpenTelemetry export.

## Requirements

- Go version from `go.mod`
- MySQL 8 for auth and migrations
- Redis only when `CACHE_STORE=redis` or running workers
- [Air](https://github.com/air-verse/air) only for live reload during development

## Quick start

Run the commands from the repository root. Copy the example, set `DB_PASS` and `MYSQL_ROOT_PASSWORD`, and generate a unique JWT secret:

```sh
cp .env.example .env
openssl rand -hex 32 # Paste the output into JWT_SECRET in .env
```

For containers, provide the same application settings as environment variables; the `.env` file is optional. Docker Compose still needs `DB_PASS` and `MYSQL_ROOT_PASSWORD` in its environment.

Start a local MySQL instance (and Redis if needed), then initialize and run the app:

```sh
docker compose --env-file .env -f deployments/docker-compose.yml up -d mysql redis
go mod download
go run . migrate
printf '%s\n' 'choose-a-strong-password' | go run . create_user -u admin -e admin@example.com --password-stdin --admin
go run . server
```

In another terminal:

```sh
curl http://localhost:8000/ping
```

The server listens on `APP_PORT` (default 8000). For local development, set `CACHE_STORE=local`; for multiple instances, set `CACHE_STORE=redis` and supply Redis settings. Run `go run . worker` when Redis is configured; the example schedule prints `schedule` every 15 seconds. If you use an existing MySQL server, create its database and user before `go run . migrate`.

### Start MySQL and Redis with Docker Compose

The quick start uses Compose to create MySQL and Redis. To start them again later, run:

```sh
docker compose --env-file .env -f deployments/docker-compose.yml up -d mysql redis
```

The app runs on your host, so keep `DB_HOST=127.0.0.1` and `DB_PORT=3306` in `.env`. Redis is optional while `CACHE_STORE=local`. To use Redis, set `CACHE_STORE=redis`, `REDIS_HOST=127.0.0.1`, and `REDIS_PORT=6379`. Stop the services with `docker compose --env-file .env -f deployments/docker-compose.yml down`.

### Live reload with Air

The project targets the Go version declared in `go.mod`. Install Air and put the Go binary directory on your `PATH`:

```sh
go install github.com/air-verse/air@v1.67.4
export PATH="$(go env GOPATH)/bin:$PATH"
```

After completing the database setup and migrations above, run this from the repository root:

```sh
air -c .air.toml
```

`.air.toml` builds the application and passes the `server` command to it. Air watches Go source files and restarts the server after changes. Press Ctrl+C to stop. If you do not need live reload, use `go run . server` instead. Air is for local development, not production deployment.

## API

- `GET /ping` — health response
- `POST /api/v1/auth/login` — access token
- `POST /api/v1/auth/register` — create a user; queue a welcome email if mail is enabled
- `POST /api/v1/auth/forgot-password` — request a reset email without disclosing account existence
- `POST /api/v1/auth/reset-password` — set a new password with the emailed token
- `GET /api/v1/auth/me` — current user
- `POST /api/v1/auth/logout` — revoke token
- `POST /api/v1/auth/change-pass` — change password
- `/api/v1/rbac/*` — users, roles and permissions
- `POST /api/v1/graphql` — authenticated GraphQL query endpoint
- `GET /api/v1/public/enum-options/:key` — public enum catalog lookup
- `GET /api/v1/graphql/entity-options` — authenticated, permission-scoped entity options
- `GET /swagger/index.html` — Swagger UI for the REST API (disabled in `release`)
- `GET /openapi.yaml` — OpenAPI 3.0 specification (disabled in `release`)

Pass `Authorization: Bearer <token>` to protected endpoints. Create roles and permissions according to your application. The base contains no privileged usernames. The `--admin` option grants an explicit admin role during trusted local setup; omit it for normal users.

The `admin` role is reserved for trusted CLI setup. API clients cannot create, rename, or delete that role; only an admin may grant or remove it. Changing a password invalidates existing access tokens for that user.

### Background jobs and email example

The scheduler example is registered in `internal/app/worker/register/schedule.go` and prints `schedule` at second 0, 15, 30, and 45 of each minute while `go run . worker` is running. It uses a named Redis lock to prevent overlapping runs across workers; schedules are not an exactly-once delivery mechanism. The example is deliberately visible; remove it when you add your own schedules.

To try the email jobs locally, start Redis and Mailpit, then enable mail in `.env`:

```sh
docker compose --env-file .env -f deployments/docker-compose.yml up -d redis mailpit
```

Set `CACHE_STORE=redis`, `REDIS_HOST=127.0.0.1`, `REDIS_PORT=6379`, `MAIL_ENABLED=true`, `MAIL_HOST=127.0.0.1`, `MAIL_PORT=1025`, `MAIL_FROM=noreply@example.test`, `MAIL_TLS_MODE=none`, and `PASSWORD_RESET_URL=http://localhost:3000/reset-password`. Run `go run . server` and `go run . worker` in separate terminals. Mailpit is available at `http://127.0.0.1:8025`; no email leaves your machine. `PASSWORD_RESET_URL` is the URL of your client application's reset page; the client reads its `token` query parameter and submits it to the API.

For example, register and request a password reset:

```sh
curl -X POST http://localhost:8000/api/v1/auth/register -H 'Content-Type: application/json' \
  -d '{"username":"alice","email":"alice@example.test","password":"StrongPass1!"}'
curl -X POST http://localhost:8000/api/v1/auth/forgot-password -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.test"}'
```

Open the reset email in Mailpit, copy the `token` query value from its link, and submit it once:

```sh
curl -X POST http://localhost:8000/api/v1/auth/reset-password -H 'Content-Type: application/json' \
  -d '{"token":"PASTE_TOKEN","new_password":"AnotherStrong1!","confirm_password":"AnotherStrong1!"}'
```

The reset token expires after 15 minutes and becomes invalid when the password changes. The queue stores only an email address for reset requests; the worker creates the token when sending. Requests for unknown addresses get the same HTTP response. If the queue is unavailable, forgot-password returns 503 with an auth error code. Registration still creates the user and returns `data.email_queued=false` if its welcome email could not be queued; check the server log for the cause. A queue delivery may be retried and can send an email more than once if SMTP accepts it before the worker sees a failure. The reset token remains usable only once. The public register and forgot-password routes each allow five requests per minute per client IP, and reset-password allows ten.

Mail jobs retry transient handler errors up to five times and then move to Asynq's archived queue for inspection. Worker logs include the task ID, type, queue, retry count, and outcome without logging the task payload. A queued job survives a worker restart while Redis retains its data. This example does not use a transactional outbox: registration can commit while enqueue fails, and a Redis data loss can remove queued jobs. Add an outbox and durable Redis configuration if your delivery requirements demand recovery across those failures. SMTP delivery can fail permanently or be duplicated: the receiving server may accept a message even if its final acknowledgment is lost.

For production SMTP, configure `MAIL_TLS_MODE=starttls` or `implicit`, use a real `MAIL_FROM`, and set an HTTPS `PASSWORD_RESET_URL`. Keep SMTP credentials in environment variables or a secret manager. If `MAIL_ENABLED=false`, registration still works without email and the forgot-password route returns 503.

### Swagger and GraphQL

In `local` or `debug` mode, open `http://localhost:8000/swagger/index.html` to explore the REST API. The UI reads the bundled [OpenAPI specification](internal/routes/openapi.yaml). Both documentation routes are unavailable in `release`. Use its **Authorize** button to send a JWT to protected routes. Registration, forgot-password, and reset-password are public operations with request and response examples; set `X-Language: en` in Swagger UI to match the English response examples. To see emails after using **Try it out**, start Redis, Mailpit, and the worker as described above. Keep the specification aligned with registered REST routes.

GraphQL uses the same access token as REST. For example:

```sh
curl -X POST http://localhost:8000/api/v1/graphql \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"query":"{ me { id username email status } }"}'
```

The [GraphQL schema](internal/app/graphql/auth/schema.graphqls) exposes `me`, `enum_keys`, `enum_options`, and `entity_options`. The public REST enum catalog currently contains `user_status`. Entity keys `users`, `roles`, and `permissions` require the existing `user:index`, `role:index`, and `permission:index` permissions respectively, both through GraphQL and the REST entity-options endpoint. GraphQL HTTP routes live in `internal/app/v1/graphql`; the in-memory catalog and access rules live in `internal/app/core/graphql`, and GORM queries live in `internal/repository/graphql`. The POST route requires JWT, permits 100 requests per minute per client IP, caps the body at 64 KiB, and gives execution a 2-second timeout. Query complexity, depth, alias count, fragment count, and parser token count are bounded; introspection is available outside release mode. To regenerate schema code, run `go run github.com/99designs/gqlgen generate --config gqlgen.yml` from the repository root.

For dropdowns, call `GET /api/v1/public/enum-options/user_status` or query `entity_options(key: "roles", page: 1, per_page: 20) { key options { id code name label extra } meta { page per_page last_page total } }`. Both entity endpoints use the same permission checks and pagination logic. The ERP Vue client in the reference export is not part of this Go repository.

## Development

```sh
make vet
make build
go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...
./build/ggb server
```

`config/` maps environment variables. `internal/app/core/auth/` owns reusable auth logic, including the mail use cases; `internal/app/core/worker/` defines task types, parameters, and dispatch. `internal/app/v1/auth/` assembles the HTTP auth module and owns its routes, handlers, and contracts. `internal/app/worker/auth/` handles background auth tasks, while `internal/app/worker/register/` registers handlers and schedules. `internal/request/` and `internal/response/` hold shared HTTP helpers; `internal/utils/` holds transport-independent parsing helpers. `internal/repository/auth/` implements auth persistence with GORM; `internal/infra/` provides resource clients and adapters. `internal/migrations/` contains auth schema migrations. The Go module path is `github.com/truongbo17/go-gin-boilerplate`.

To add an HTTP feature, create its `NewModule` and `RegisterRoutes` methods under `internal/app/v1/<feature>`, then construct and mount it in `internal/routes/main.go`. Keep repositories, services, controllers, and route middleware inside the feature. See [Architecture and extension points](docs/architecture.md#extend-the-application) for worker setup too.

See [Architecture and extension points](docs/architecture.md) for package boundaries, feature workflows, optional infrastructure, and the local performance reference. See [Go conventions](docs/go-style.md) for naming, context, error, and review guidelines.

### Infrastructure included

| Package | Purpose |
| --- | --- |
| `internal/infra/database` | MySQL and optional PostgreSQL connection pools with tracing |
| `internal/infra/mongodb`, `clickhouse` | Optional document and analytics database clients |
| `internal/infra/kafka` | Optional Kafka producer, consumer, and connectivity check |
| `internal/infra/objectstore` | Optional S3 or compatible object storage client |
| `internal/infra/cache`, `redis` | Local or shared cache and Redis connection |
| `internal/middlewares/limiter` | Per-route request limits backed by memory or Redis |
| `internal/infra/http` | Outbound HTTP client with timeout, tracing, and an 8 MiB response limit |
| `internal/infra/worker`, `schedule` | Optional Redis-backed jobs and schedules |
| `internal/infra/mail` | SMTP sender for optional auth email jobs |
| `internal/infra/logger`, `tracer` | Logging and telemetry |

Add feature-specific persistence under `internal/repository/<feature>/` and external-system integrations under `internal/infra/`; keep shared HTTP helpers under `internal/request/` and `internal/response/`, common language messages under `internal/i18n/`, and business rules in `internal/app/core/`. Feature error messages and their HTTP status mapping live with that feature's responses. `/ping` checks the HTTP process.

See [Optional infrastructure clients](docs/architecture.md#optional-infrastructure) for constructors and shutdown. These clients are available for features that need them; the default server does not connect to all backends. The bundled auth schema migrations are MySQL-specific, so PostgreSQL support here is an independent connection client rather than a drop-in auth database switch.

The MySQL pool defaults to 30 open and 15 idle connections. Adjust `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME_MINUTES`, and `DB_CONN_MAX_IDLE_TIME_MINUTES` for your database limits and workload. Set an idle or lifetime value to `0` to disable that limit. The worker stops schedules before draining in-flight tasks on SIGINT/SIGTERM; task shutdown waits up to 25 seconds.

When tracing is enabled, `TRACER_SAMPLE_RATIO` controls the share of new traces sampled (default `0.1`).

The Docker Compose file under `deployments/` starts MySQL and Redis for development. Copy `.env.example` to `.env` and set secrets before using it. GitHub Actions checks formatting, source, vet, build, dependencies, and CodeQL; Dependabot opens dependency update PRs and security alerts are enabled.

## Security defaults

- Supply a unique JWT secret with at least 32 characters; no secret is bundled.
- Set exact CORS origins for your clients. The example allows `http://localhost:3000`.
- Use Redis-backed cache for token revocation and rate limits across replicas.
- Login is limited to 10 attempts per minute per client IP; each other auth or RBAC endpoint has its own 300 requests per minute per IP limit, declared alongside that route in `internal/app/v1/auth/routes.go`. `/ping` is not rate limited so health probes cannot exhaust client quotas. A 429 response includes `Retry-After`; rate limit and request ID headers are exposed to allowed browser origins.
- By default, forwarded IP headers are ignored. Behind a reverse proxy, set `APP_TRUSTED_PROXIES` to its exact IP or CIDR (comma-separated for multiple proxies), and restrict direct access to the app. Never set it to a public or unrestricted CIDR: the rate limiter uses the resulting client IP. Configure HTTPS and HSTS at the TLS-terminating proxy.
- Request bodies are limited to 1 MiB, including streamed bodies. API failures use HTTP error status codes, and responses use `Cache-Control: no-store`.
- Access tokens include a password version; changing a password invalidates tokens issued before the change. Tokens issued by older versions of this boilerplate need a fresh login after upgrading.
- Place credentials only in local environment or a secret manager, never in Git.
- Use `--password-stdin` when creating users so passwords do not appear in process arguments or shell history.
