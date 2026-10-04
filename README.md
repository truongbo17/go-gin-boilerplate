# Go Gin Boilerplate

A reusable Gin API base derived from a production Go service. It includes JWT authentication, role-based access control, MySQL migrations, optional Redis cache and workers, rate limiting, structured logs, request IDs, and optional OpenTelemetry export.

## Requirements

- Go version from `go.mod`
- MySQL 8 for auth and migrations
- Redis only when `CACHE_STORE=redis` or running workers

## Quick start

```sh
cp .env.example .env
# Set DB_USER, DB_PASS, DB_DATABASE and a unique 32+ character JWT_SECRET in .env.
go mod download
go run . migrate
go run . create_user -u admin -e admin@example.com -p 'choose-a-strong-password' --admin
go run . server
curl http://localhost:8000/ping
```

Create the MySQL database and user before running migrations. The server listens on `APP_PORT` (default 8000). For local development, set `CACHE_STORE=local`; for multiple instances, set `CACHE_STORE=redis` and supply Redis settings. Run `go run . worker` only when Redis is configured and your application has registered jobs in `internal/app/core/register`.

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
```

`config/` maps environment variables. `internal/app/core/auth/` owns reusable auth logic; `internal/app/v1/auth/` owns HTTP routes. `internal/infra/` contains adapters. `internal/migrations/` contains auth schema migrations. The Go module path is `github.com/truongbo17/go-gin-boilerplate`.

The Docker Compose file under `deployments/` starts MySQL and Redis for development. Copy `.env.example` to `.env` and set secrets before using it. GitHub Actions runs tests, vet, build, a public-source check, and CodeQL; Dependabot opens dependency update PRs.

## Security defaults

- Supply a unique JWT secret with at least 32 characters; no secret is bundled.
- Set exact CORS origins for your clients. The example allows `http://localhost:3000`.
- Use Redis-backed cache for token revocation and rate limits across replicas.
- Place credentials only in local environment or a secret manager, never in Git.

See [base conversion criteria](docs/BASE_SPEC.md) for the retained scope.
