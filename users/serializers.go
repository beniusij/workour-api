package users

import (
	"github.com/gin-gonic/gin"
	"workour-api/common"
)

type UserSerializer struct {
	c *gin.Context
}

type UserResponse struct {
	Email		string `json:"email"`
	FirstName	string `json:"first_name"`
	LastName	string `json:"last_name"`
	Token		string `json:"token"`
}

func (u *UserSerializer) Response() UserResponse {
	userModel := u.c.MustGet("my_user_model").(User)
	user := UserResponse{
		Email:     userModel.Email,
		FirstName: userModel.FirstName,
		LastName:  userModel.LastName,
		Token:     common.GetToken(userModel.ID),
	}
	return user
}