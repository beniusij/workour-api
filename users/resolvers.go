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

	userId, err := userStruct.Save(userValidator.UserModel)

	if  err != nil {
		return nil, err
	}

	user := User{ID: userId}

	c.Set("status", http.StatusCreated)

	return user, nil
}

// GetUserResolver resolves our user query through a db call to GetById
func GetUserResolver(p g.ResolveParams) (interface{}, error) {
	user := User{}
	id := p.Args["id"].(uint)

	user, err := user.GetById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// The resolver evaluates if user with provided email exists, and if so
// authenticates with provided password. If authentication succeeds a token
// is created and set in the header, and user non-sensitive details are
// returned in the response body.
// TODO remove this resolver and related tests
func AuthenticateUserResolver(p g.ResolveParams) (interface{}, error) {
	// Get user by email
	//email := p.Args["email"].(string)
	//user, err := GetByEmail(email)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Check password
	//psw := p.Args["password"].(string)
	//err = user.CheckPassword(psw)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Generate token
	//authToken := authentication.AuthToken{}
	//token, err := authToken.GenerateToken(email)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Set up Session and return it
	//var session Session
	//session.Email = user.Email
	//session.Token = token
	//
	//c := p.Context.(*gin.Context)
	//c.Set("status", http.StatusOK)
	//
	//return session, nil
	return nil, nil
}
