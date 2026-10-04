# Production-derived Go Gin base

## Objective

Turn the local production application into a reusable, public Gin API starter. Preserve reusable application infrastructure and auth/RBAC. Remove the email marketing, SMS, contact, campaign, tracking, gallery and template products, all company identifiers, and generated API docs that describe those products.

## Commands

- Install: `go mod download`
- Test: `go test ./...`
- Static analysis: `go vet ./...`
- Build: `go build ./...`
- Run: `go run . server`

## Structure

- `cmd`: CLI entry points
- `config`: environment configuration
- `internal/app/core/auth`: domain models, repositories and services
- `internal/app/v1/auth`: HTTP auth and RBAC routes
- `internal/infra`: reusable adapters
- `internal/migrations`: base auth schema
- `internal/middlewares`, `internal/routes`: HTTP composition

## Style

Use `gofmt`, standard Go errors, and explicit dependency boundaries. Keep example configuration free of real credentials. New handlers should validate input at the HTTP boundary.

## Testing

Use Go tests for reusable logic and HTTP behavior. CI runs tests, vet, build, and a source scan for company-specific leftovers and private keys.

## Boundaries

- Always: keep a clean build, document required services, preserve no production secrets.
- Ask first: destructive changes to the original GitLab repository or local working tree.
- Never: publish real credentials, company infrastructure endpoints, customer data, or business workflows.

## Success criteria

1. `go test ./...`, `go vet ./...`, and `go build ./...` pass.
2. No company-specific names, endpoints, email marketing or SMS modules remain in the published tree.
3. The README and example environment can bootstrap the base.
4. GitHub Actions checks and Dependabot are configured.
5. Changes are split into reviewable commits on a GitHub branch.
