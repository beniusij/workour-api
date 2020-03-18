package users

import (
	"github.com/gin-gonic/gin"
	g "github.com/graphql-go/graphql"
	"net/http"
	//"workour-api/authentication"
)

type Session struct {
	Email string
	Token string
}

// Handles mutation to create a user
func CreateUserResolver(p g.ResolveParams) (interface{}, error) {
 	userValidator := NewUserValidator()
	userStruct := &User{}
	c := p.Context.(*gin.Context)

	if err := userValidator.ValidateForm(p.Args); err != nil {
		return nil, err
	}

	user, err := userStruct.Save(userValidator.UserModel)

	if  err != nil {
		return nil, err
	}

	c.Set("status", http.StatusCreated)

	return user, nil
}

// GetUserResolver resolves our user query through a db call to GetById
func GetUserResolver(p g.ResolveParams) (interface{}, error) {
	user := User{}
	user.ID = p.Args["id"].(uint)

	err := user.GetById()
	if err != nil {
		return nil, err
	}

	return user, nil
}
