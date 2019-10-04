package users

import (
	g "github.com/graphql-go/graphql"
)

// Handles mutation to create a user
func CreateUserResolver() func(p g.ResolveParams) (interface{}, error) {
	return func(p g.ResolveParams) (interface{}, error) {
		//fmt.Printf("%v", reflect.TypeOf(p.Args))
		validatedUser, err := validateUserForm(p.Args)
		if err != nil {
			return nil, err
		}

		user := &User{
			Email: validatedUser.Email,
		}
		return nil, nil
	}
}

// GetUserResolver resolves our user query through a db call to GetUserById
func GetUserResolver() func(p g.ResolveParams) (interface{}, error) {
	return func(p g.ResolveParams) (interface{}, error) {
		// Strip the name from arguments and assert that it is an int
		id, ok := p.Args["id"].(int)

		if ok {
			user, err := GetUserById(id)
			return user, err
		}

		return nil, nil
	}
}
