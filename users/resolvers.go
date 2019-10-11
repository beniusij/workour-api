package users

import (
	"github.com/gin-gonic/gin"
	g "github.com/graphql-go/graphql"
	"net/http"
)

// Handles mutation to create a user
func CreateUserResolver(p g.ResolveParams) (interface{}, error) {
	userValidator := NewUserValidator()
	user := &User{}
	c := p.Context.(*gin.Context)

	if err := userValidator.ValidateForm(p.Args); err != nil {
		return nil, err
	}

	if _, err := user.SaveEntity(&userValidator.UserModel); err != nil {
		return nil, err
	}

	c.Set("status", http.StatusCreated)

	return userValidator.UserModel, nil
}

// GetUserResolver resolves our user query through a db call to GetEntityById
func GetUserResolver(p g.ResolveParams) (interface{}, error) {
	user := &User{}
	id := p.Args["id"].(int)

	user, err := user.GetEntityById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}