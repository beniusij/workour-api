package users

import (
	g "github.com/graphql-go/graphql"
)

// Handles mutation to create a user
func CreateUserResolver(p g.ResolveParams) (interface{}, error) {
	userValidator := NewUserValidator()
	user := &User{}

	if err := userValidator.ValidateForm(p.Args); err != nil {
		return nil, err
	}

	if _, err := user.SaveEntity(&userValidator.UserModel); err != nil {
		return nil, err
	}

	return userValidator.UserModel.ID, nil
}

// GetUserResolver resolves our user query through a db call to GetEntityById
func GetUserResolver(p g.ResolveParams) (interface{}, error) {
	user := &User{}
	id, ok := p.Args["id"].(int)

	if ok {
		user, err := user.GetEntityById(id)
		return user, err
	}

	return nil, nil
}
