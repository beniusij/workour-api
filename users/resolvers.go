package users

import (
	g "github.com/graphql-go/graphql"
)

// Handles mutation to create a user
func CreateUserResolver(p g.ResolveParams) (interface{}, error) {
	userValidator := NewUserValidator()

	if err := userValidator.ValidateForm(p.Args); err != nil {
		return nil, err
	}

	if err := SaveUser(&userValidator.UserModel); err != nil {
		return nil, err
	}

	return userValidator.UserModel.ID, nil
}

// GetUserResolver resolves our user query through a db call to GetUserById
func GetUserResolver(p g.ResolveParams) (interface{}, error) {
	// Strip the name from arguments and assert that it is an int
	id, ok := p.Args["id"].(int)

	if ok {
		user, err := GetUserById(id)
		return user, err
	}

	return nil, nil
}
