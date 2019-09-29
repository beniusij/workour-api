package gql

import (
	g "github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	u "workour-api/users"
)

// Holds a connection to our database (???)
type Resolver struct {
	db *gorm.DB
}

// UserResolver resolves our user query through a db call to GetUserById
func (r *Resolver) UserResolver(p g.ResolveParams) (interface{}, error) {
	// Strip the name from arguments and assert that it is an int
	id, ok := p.Args["id"].(uint)
	if ok {
		user, err := u.GetUserById(id)
		return user, err
	}

	return nil, nil
}
