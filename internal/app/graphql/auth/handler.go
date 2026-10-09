package auth

import (
	"context"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/graphql/auth/generated"
	"github.com/vektah/gqlparser/v2/ast"
)

type userKey struct{}

// WithUser passes the user verified by the shared JWT middleware to resolvers.
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

// NewHandler serves only POST requests. Introspection is available for local
// development and disabled in release deployments.
func NewHandler(env string) http.Handler {
	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{}}))
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.SetParserTokenLimit(10000)
	srv.Use(extension.FixedComplexityLimit(20))
	if env != config.ReleaseMode {
		srv.Use(extension.Introspection{})
	}
	return srv
}
