package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/sirupsen/logrus"
	"github.com/truongbo17/go-gin-boilerplate/config"
	"github.com/truongbo17/go-gin-boilerplate/internal/app/core/auth/models"
	coregraphql "github.com/truongbo17/go-gin-boilerplate/internal/app/core/graphql"
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
func NewHandler(env string, metadata coregraphql.Service, log *logrus.Logger) http.Handler {
	schema := generated.Config{Resolvers: &Resolver{Metadata: metadata, Logger: log}}
	schema.Complexity.Query.EntityOptions = func(childComplexity int, _ string, _ *string, _ *int, perPage *int) int {
		limit := 20
		if perPage != nil {
			limit = *perPage
		}
		return 10 + limit/10 + childComplexity
	}
	srv := handler.New(generated.NewExecutableSchema(schema))
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](200))
	srv.SetParserTokenLimit(2000)
	srv.Use(queryLimits{})
	srv.Use(extension.FixedComplexityLimit(100))
	if env != config.ReleaseMode {
		srv.Use(extension.Introspection{})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		srv.ServeHTTP(w, r.WithContext(ctx))
	})
}
