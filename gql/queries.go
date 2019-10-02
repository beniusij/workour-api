package gql

import (
	g "github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/jinzhu/gorm"
)

type Root struct {
	Query *g.Object
}

// NewRoot returns base query type. This is where we add all the base queries
func NewRoot(db *gorm.DB) *Root {
	// Create a resolver holding our database. Resolver can be found in resolvers.go
	resolver := Resolver{db: db}

	// Create a new Root that describes our base query set up. In this
	// example we have a user query that takes one argument called ID
	root := Root{
		Query: g.NewObject(
			g.ObjectConfig{
				Name:        "Query",
				Fields:      g.Fields{
					"user": &g.Field{
						// Slice of User type which can be found in types.go
						Type: User,
						Args: g.FieldConfigArgument{
							"id": &g.ArgumentConfig{
								Type:         g.Int,
							},
						},
						Resolve: resolver.UserResolver,
					},
				},
			},
		),
	}

	return &root
}

// This one runs our graphql queries
func ExecuteQuery(query string, schema g.Schema) (*g.Result, gqlerrors.FormattedErrors) {
	result := g.Do(g.Params{
		Schema:			schema,
		RequestString:	query,
	})

	return result, result.Errors
}