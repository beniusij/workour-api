package gql

import g "github.com/graphql-go/graphql"

var User = g.NewObject(
	g.ObjectConfig{
		Name:        "User",
		Fields:      g.Fields{
			"id": &g.Field{
				Type: g.Int,
			},
			"email": &g.Field{
				Type: g.String,
			},
			"first_name": &g.Field{
				Type: g.String,
			},
			"last_name": &g.Field{
				Type: g.String,
			},
		},
	},
)