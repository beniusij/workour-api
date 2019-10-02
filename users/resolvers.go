package users

import (
	"github.com/graphql-go/graphql"
)

// UserResolver resolves our user query through a db call to GetUserById
func UserResolver() func(p graphql.ResolveParams) (interface{}, error) {
	return func(p graphql.ResolveParams) (interface{}, error) {
	// Strip the name from arguments and assert that it is an int
	id, ok := p.Args["id"].(int)

	if ok {
	user, err := GetUserById(id)
	return user, err
	}

	return nil, nil
	}
}
