package gql

import g "github.com/graphql-go/graphql"

var UserType = g.NewObject(
	g.ObjectConfig{
		Name:        "User",
		Fields:      g.Fields{
			"ID": &g.Field{
				Type: g.Int,
			},
			"Email": &g.Field{
				Type: g.String,
			},
			"FirstName": &g.Field{
				Type: g.String,
			},
			"LastName": &g.Field{
				Type: g.String,
			},
			"Token": &g.Field{
				Type: g.String,
			},
		},
	},
)

var SessionType = g.NewObject(
	g.ObjectConfig{
		Name:	"Session",
		Fields:	g.Fields{
			"Token": &g.Field{
				Type: g.String,
			},
			"Email": &g.Field{
				Type: g.String,
			},
		},
	},
)