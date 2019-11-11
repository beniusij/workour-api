package gql

import (
	"github.com/gin-gonic/gin"
	g "github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	u "workour-api/users"
)

type Root struct {
	Query *g.Object
	Mutation *g.Object
}

var queryType = g.NewObject(
	g.ObjectConfig{
		Name: "Query",
		Fields: g.Fields{
			"user": &g.Field{
				Type: UserType,
				Args: g.FieldConfigArgument{
					"id": &g.ArgumentConfig{
						Type: g.Int,
					},
				},
				Resolve:     u.GetUserResolver,
				Description: "Get user by id",
			},
		},
	},
)

var mutationType = g.NewObject(
	g.ObjectConfig{
		Name: "Mutation",
		Fields: g.Fields{
			"register": &g.Field{
				Type: UserType,
				Args: g.FieldConfigArgument{
					"email": &g.ArgumentConfig{
						Type: g.NewNonNull(g.String),
					},
					"first_name": &g.ArgumentConfig{
						Type: g.NewNonNull(g.String),
					},
					"last_name": &g.ArgumentConfig{
						Type: g.NewNonNull(g.String),
					},
					"password": &g.ArgumentConfig{
						Type: g.NewNonNull(g.String),
					},
					"password_confirm": &g.ArgumentConfig{
						Type: g.NewNonNull(g.String),
					},
				},
				Resolve: u.CreateUserResolver,
				Description: "Create a new user",
			},
			"login": &g.Field{
				Name: 			"Login",
				Type: 			UserType,
				Args: 			g.FieldConfigArgument{
					"email": 	&g.ArgumentConfig{
						Type: 			g.NewNonNull(g.String),
						Description:	"User email",
					},
					"password":	&g.ArgumentConfig{
						Type:         	g.NewNonNull(g.String),
						Description:  	"User password in plain text",
					},
				},
				Resolve: 		u.AuthenticateUserResolver,
				Description: 	"A login used by the user-facing site to create user JSON Web Token",
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
		Mutation: mutationType,
	}

	return &root
}

// This one runs our graphql queries
func ExecuteQuery(query string, v map[string]interface{}, schema g.Schema, c *gin.Context) (*g.Result, gqlerrors.FormattedErrors) {
	result := g.Do(g.Params{
		Context:		c,
		Schema:			schema,
		RequestString:	query,
		VariableValues:	v,
	})

	return result, result.Errors
}