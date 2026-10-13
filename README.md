# Go Gin Boilerplate

A reusable Go boilerplate built with Gin. It includes REST and GraphQL APIs, JWT authentication, role-based access
control (RBAC), MySQL migrations, rate limiting, and background workers. Features are organized so you can replace the
sample authentication module or add your own business logic.

## Quick start

### Requirements

- Go version declared in `go.mod`
- Docker Compose for MySQL; Redis is needed only for workers or when Redis is selected as the cache store
- Python 3 to enter the admin password without displaying it in the terminal

### 1. Configure your local environment

```sh
cp .env.example .env
openssl rand -hex 32
```

Put the generated value in `JWT_SECRET` in `.env`. Set values for `DB_PASS` and `MYSQL_ROOT_PASSWORD` as well.

### 2. Start MySQL

```sh
docker compose --env-file .env -f deployments/docker-compose.yml up -d mysql
```

The application runs on your host, so keep `DB_HOST=127.0.0.1` and `DB_PORT=3306` in `.env`.

### 3. Create the database tables and an admin account

```sh
go run . migrate
python3 -c 'import getpass; print(getpass.getpass("Admin password: "))' | go run . create_user \
  --username admin \
  --email admin@example.com \
  --password-stdin \
  --admin
```

Use `--admin` only for a trusted operator. The API does not allow clients to create, rename, or delete the `admin` role.

### 4. Run the API

```sh
go run . server
```

The API listens on `http://localhost:8000` by default. Check that it is running:

```sh
curl http://localhost:8000/ping
```

When `APP_ENV` is `local` or `debug`, open `http://localhost:8000/swagger/index.html` to explore the REST API. Swagger
and OpenAPI endpoints are disabled in `release` mode.

## Included endpoints

| Method and path                        | Description                           |
|----------------------------------------|---------------------------------------|
| `POST /api/v1/auth/register`           | Create an account                     |
| `POST /api/v1/auth/login`              | Sign in and receive a JWT             |
| `POST /api/v1/auth/forgot-password`    | Request a password reset              |
| `POST /api/v1/auth/reset-password`     | Set a new password with a reset token |
| `/api/v1/rbac/*`                       | Manage users, roles, and permissions  |
| `POST /api/v1/graphql`                 | Authenticated GraphQL endpoint        |
| `GET /api/v1/public/enum-options/:key` | Read a public enum catalog entry      |
| `GET /api/v1/graphql/entity-options`   | List entities with permission checks  |

Send `Authorization: Bearer <token>` to protected endpoints. Changing a password invalidates the user's existing access
tokens.

## Workers and sample email jobs

Workers require Redis. Set `CACHE_STORE=redis`, `REDIS_HOST=127.0.0.1`, and `REDIS_PORT=6379` in `.env`, then run:

```sh
docker compose --env-file .env -f deployments/docker-compose.yml up -d redis
go run . worker
```

The worker includes a sample schedule that prints `schedule` every 15 seconds, plus welcome and password-reset email
jobs.

To test email locally, start Mailpit:

```sh
docker compose --env-file .env -f deployments/docker-compose.yml up -d mailpit
```

Set `MAIL_ENABLED=true`, `MAIL_HOST=127.0.0.1`, `MAIL_PORT=1025`, `MAIL_FROM=noreply@example.test`,
`MAIL_TLS_MODE=none`, and `PASSWORD_RESET_URL=http://localhost:3000/reset-password`. Run `server` and `worker` in
separate terminals. View messages at `http://127.0.0.1:8025`; Mailpit keeps them local.

For production, use `MAIL_TLS_MODE=starttls` or `implicit`, an HTTPS reset URL, and a secret manager for SMTP
credentials. The queue retries jobs but does not guarantee exactly-once delivery; an email may be sent more than once if
the SMTP connection fails while the result is being acknowledged.

## Common configuration

Settings can be provided through `.env` or environment variables. Common options include:

- `APP_ENV`: `local`, `debug`, or `release`
- `APP_PORT`: HTTP port; defaults to `8000`
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_DATABASE`: MySQL connection settings
- `CACHE_STORE`: `local` for a single instance, or `redis` to share cache and rate-limit state across instances
- `JWT_SECRET`: a private signing key of at least 32 characters
- `MAIL_ENABLED`: enable or disable email jobs
- `TRACER_ENABLE`: enable OpenTelemetry export

See [.env.example](.env.example) for the full list. Do not commit `.env` files or real secrets.

## Source layout

- `cmd/`: `server`, `worker`, `migrate`, `create_user`, and `version` commands and resource lifecycle
- `config/`: configuration loading and validation
- `internal/app/core/<feature>/`: business logic and types independent of HTTP
- `internal/app/v1/<feature>/`: REST routes, controllers, requests, and responses
- `internal/app/worker/<feature>/`: job handlers that call business services
- `internal/repository/<feature>/`: database queries and transactions
- `internal/infra/`: clients for databases, Redis, Kafka, MongoDB, ClickHouse, S3, SMTP, tracing, and workers
- `internal/routes/`: Gin setup, shared middleware, and HTTP feature composition

To add an HTTP feature, create its module and `RegisterRoutes` under `internal/app/v1/<feature>`, then initialize and
mount it in `internal/routes/main.go`. Keep rate limits, authentication, permission checks, validation, and controller
calls next to each route so the request flow is easy to follow.

See [docs/architecture.md](docs/architecture.md) for request flow, workers, and infrastructure.
See [docs/go-style.md](docs/go-style.md) for coding conventions.

## Development commands

```sh
make fmt
make vet
make build
go run . version
```

Air can rebuild and restart the server when source files change:

```sh
go install github.com/air-verse/air@v1.67.4
air -c .air.toml
```

## Security defaults

- Use a unique `JWT_SECRET` for each environment; do not use an example key.
- Allow only the CORS origins required by your frontend.
- Behind a reverse proxy, trust only its specific IP addresses or CIDR ranges in `APP_TRUSTED_PROXIES`.
- Use Redis-backed cache when running multiple instances so token revocation and rate limits are shared.
- Terminate TLS at the reverse proxy. Do not expose the application port directly to the Internet when the proxy
  controls external access.
- Create accounts through the CLI. Use `--password-stdin` so passwords do not appear in shell history.
