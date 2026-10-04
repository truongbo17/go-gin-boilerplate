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

Start a local MySQL instance (and Redis if needed), then initialize and run the app:

```sh
docker compose --env-file .env -f deployments/docker-compose.yml up -d mysql redis
go mod download
go run . migrate
go run . create_user -u admin -e admin@example.com -p 'choose-a-strong-password' --admin
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

The project targets Go 1.23.4. Install a compatible Air release and put the Go binary directory on your `PATH`:

```sh
go install github.com/air-verse/air@v1.61.7
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
- `GET /api/v1/auth/me` — current user
- `POST /api/v1/auth/logout` — revoke token
- `POST /api/v1/auth/change-pass` — change password
- `/api/v1/rbac/*` — users, roles and permissions

Pass `Authorization: Bearer <token>` to protected endpoints. Create roles and permissions according to your application. The base contains no privileged usernames. The `--admin` option grants an explicit admin role during trusted local setup; omit it for normal users.

## Development

```sh
make test
make vet
make build
./build/ggb server
```

`config/` maps environment variables. `internal/app/core/auth/` owns reusable auth logic; `internal/app/v1/auth/` owns HTTP routes. `internal/infra/` contains adapters. `internal/migrations/` contains auth schema migrations. The Go module path is `github.com/truongbo17/go-gin-boilerplate`.

The Docker Compose file under `deployments/` starts MySQL and Redis for development. Copy `.env.example` to `.env` and set secrets before using it. GitHub Actions runs tests, vet, build, a public-source check, and CodeQL; Dependabot opens dependency update PRs.

## Security defaults

- Supply a unique JWT secret with at least 32 characters; no secret is bundled.
- Set exact CORS origins for your clients. The example allows `http://localhost:3000`.
- Use Redis-backed cache for token revocation and rate limits across replicas.
- Place credentials only in local environment or a secret manager, never in Git.
