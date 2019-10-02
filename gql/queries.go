package gql

import (
	g "github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	u "workour-api/users"
)

type Root struct {
	Query *g.Object
}

var queryType = g.NewObject(
	g.ObjectConfig{
		Name: "Query",
		Fields: g.Fields{
			"user": &g.Field{
				Type: u.UserType,
				Args: g.FieldConfigArgument{
					"id": &g.ArgumentConfig{
						Type: g.Int,
					},
				},
				Resolve:     u.UserResolver(),
				Description: "Get user by id",
			},
		},
	},
)

// NewRoot returns base query type. This is where we add all the base queries
func NewRoot() *Root {
	// Create a new Root that describes our base query set up. In this
	// example we have a user query that takes one argument called ID
	root := Root{
		Query: queryType,
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