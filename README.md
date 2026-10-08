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

The server listens on `APP_PORT` (default 8000). For local development, set `CACHE_STORE=local`; for multiple instances, set `CACHE_STORE=redis` and supply Redis settings. Run `go run . worker` only when Redis is configured and your application has registered jobs in `internal/app/core/register`. If you use an existing MySQL server, create its database and user before `go run . migrate`.

### Start MySQL and Redis with Docker Compose

The quick start uses Compose to create MySQL and Redis. To start them again later, run:

```sh
docker compose --env-file .env -f deployments/docker-compose.yml up -d mysql redis
```

The app runs on your host, so keep `DB_HOST=127.0.0.1` and `DB_PORT=3306` in `.env`. Redis is optional while `CACHE_STORE=local`. To use Redis, set `CACHE_STORE=redis`, `REDIS_HOST=127.0.0.1`, and `REDIS_PORT=6379`. Stop the services with `docker compose --env-file .env -f deployments/docker-compose.yml down`.

### Live reload with Air

The project targets Go 1.26.8. Install Air and put the Go binary directory on your `PATH`:

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
- `GET /ready` — dependency readiness (MySQL and Redis when enabled)
- `POST /api/v1/auth/login` — access token
- `GET /api/v1/auth/me` — current user
- `POST /api/v1/auth/logout` — revoke token
- `POST /api/v1/auth/change-pass` — change password
- `/api/v1/rbac/*` — users, roles and permissions

Pass `Authorization: Bearer <token>` to protected endpoints. Create roles and permissions according to your application. The base contains no privileged usernames. The `--admin` option grants an explicit admin role during trusted local setup; omit it for normal users.

The `admin` role is reserved for trusted CLI setup. API clients cannot create, rename, or delete that role; only an admin may grant or remove it. Changing a password invalidates existing access tokens for that user.

## Development

```sh
make test
make vet
make build
go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...
./build/ggb server
```

`config/` maps environment variables. `internal/app/core/auth/` owns reusable auth logic; `internal/app/v1/auth/` owns HTTP routes. `internal/infra/` contains adapters. `internal/migrations/` contains auth schema migrations. The Go module path is `github.com/truongbo17/go-gin-boilerplate`.

See [Go conventions](docs/go-style.md) for the project's naming, context, error, and package-boundary guidelines.

### Infrastructure included

| Package | Purpose |
| --- | --- |
| `internal/infra/database` | MySQL connection pool and tracing |
| `internal/infra/cache`, `redis`, `limiter` | Local or shared cache, Redis, and request limits |
| `internal/infra/health` | Readiness check for configured dependencies |
| `internal/infra/http` | Outbound HTTP client with timeout, tracing, and an 8 MiB response limit |
| `internal/infra/worker`, `schedule` | Optional Redis-backed jobs and schedules |
| `internal/infra/logger`, `tracer`, `i18n` | Logging, telemetry, and messages |

Add application-specific adapters under `internal/infra/` when an actual integration needs them; keep business rules in `internal/app/core/`. `/ping` checks the HTTP process, while `/ready` returns 503 if MySQL or configured Redis is unavailable.

The MySQL pool defaults to 30 open and 15 idle connections. Adjust `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME_MINUTES`, and `DB_CONN_MAX_IDLE_TIME_MINUTES` for your database limits and workload. Set an idle or lifetime value to `0` to disable that limit. The worker stops schedules before draining in-flight tasks on SIGINT/SIGTERM; task shutdown waits up to 25 seconds.

When tracing is enabled, `TRACER_SAMPLE_RATIO` controls the share of new traces sampled (default `0.1`).

The Docker Compose file under `deployments/` starts MySQL and Redis for development. Copy `.env.example` to `.env` and set secrets before using it. GitHub Actions runs tests, vet, build, a source check, govulncheck, and CodeQL; Dependabot opens dependency update PRs and security alerts are enabled.

## Security defaults

- Supply a unique JWT secret with at least 32 characters; no secret is bundled.
- Set exact CORS origins for your clients. The example allows `http://localhost:3000`.
- Use Redis-backed cache for token revocation and rate limits across replicas.
- Login is limited to 10 attempts per minute per client IP; authenticated API routes use a separate 300 requests per minute per IP limit. `/ping` and `/ready` are not rate limited so health probes cannot exhaust client quotas. A 429 response includes `Retry-After`; rate limit and request ID headers are exposed to allowed browser origins.
- By default, forwarded IP headers are ignored. Behind a reverse proxy, set `APP_TRUSTED_PROXIES` to its exact IP or CIDR (comma-separated for multiple proxies), and restrict direct access to the app. Never set it to a public or unrestricted CIDR: the rate limiter uses the resulting client IP. Configure HTTPS and HSTS at the TLS-terminating proxy.
- Request bodies are limited to 1 MiB, including streamed bodies. API failures use HTTP error status codes, and responses use `Cache-Control: no-store`.
- Access tokens include a password version; changing a password invalidates tokens issued before the change. Tokens issued by older versions of this boilerplate need a fresh login after upgrading.
- Place credentials only in local environment or a secret manager, never in Git.
- Use `--password-stdin` when creating users so passwords do not appear in process arguments or shell history.
