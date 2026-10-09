package auth

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// queryLimits runs on gqlgen's validated AST before any resolver executes.
type queryLimits struct{}

func (queryLimits) ExtensionName() string                   { return "QueryLimits" }
func (queryLimits) Validate(graphql.ExecutableSchema) error { return nil }

func (queryLimits) MutateOperationContext(_ context.Context, op *graphql.OperationContext) *gqlerror.Error {
	if len(op.Doc.Fragments) > 50 {
		return gqlerror.Errorf("too many fragments")
	}
	operation := op.Doc.Operations.ForName(op.OperationName)
	if operation == nil {
		return nil // gqlgen reports the invalid operation name.
	}
	fragments := make(map[string]ast.SelectionSet, len(op.Doc.Fragments))
	for _, fragment := range op.Doc.Fragments {
		fragments[fragment.Name] = fragment.SelectionSet
	}
	aliases := 0
	active := make(map[string]bool)
	var walk func(ast.SelectionSet, int) *gqlerror.Error
	walk = func(selections ast.SelectionSet, depth int) *gqlerror.Error {
		if len(selections) == 0 {
			return nil
		}
		if depth > 8 {
			return gqlerror.Errorf("query depth exceeds 8")
		}
		for _, selection := range selections {
			switch field := selection.(type) {
			case *ast.Field:
				if field.Alias != "" && field.Alias != field.Name {
					aliases++
					if aliases > 100 {
						return gqlerror.Errorf("too many aliases")
					}
				}
				if err := walk(field.SelectionSet, depth+1); err != nil {
					return err
				}
			case *ast.InlineFragment:
				if err := walk(field.SelectionSet, depth); err != nil {
					return err
				}
			case *ast.FragmentSpread:
				if active[field.Name] {
					return gqlerror.Errorf("fragment cycle")
				}
				active[field.Name] = true
				if err := walk(fragments[field.Name], depth); err != nil {
					return err
				}
				delete(active, field.Name)
			}
		}
		return nil
	}
	return walk(operation.SelectionSet, 1)
}
