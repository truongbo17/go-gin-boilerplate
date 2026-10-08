# Go conventions in this repository

These conventions follow [Effective Go](https://go.dev/doc/effective_go) and the [Go code review comments](https://go.dev/wiki/CodeReviewComments).

- Run `gofmt` and `go vet`; CI checks both. Keep package names short and lower case. Use `ID`, `URL`, and `HTTP` in new identifiers.
- Pass `context.Context` as the first argument to request-scoped operations. Pass the request context through database, cache, queue, and outbound HTTP calls.
- Return errors with enough operation context using `%w`. Close opened resources on startup failure and shut them down in reverse dependency order.
- Construct cheap services and repositories directly. Keep interfaces near consumers when a second implementation or test seam needs one; avoid package-wide singletons for request-scoped dependencies.
- Keep core business rules under `internal/app/core/`, HTTP handlers under `internal/app/v1/`, and external adapters under `internal/infra/`. Add an adapter when a real integration needs it.
- For changes in behavior, test the public result and important failure paths. Run `go test ./...` before committing.

Go's [`database/sql` connection pool guide](https://go.dev/doc/database/manage-connections) explains the pool settings exposed in `.env.example`. Tune them from observed database capacity and wait time rather than increasing them automatically.
