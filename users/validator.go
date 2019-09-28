package users

import (
	"errors"
	"github.com/gin-gonic/gin"
	"workour-api/common"
)

type UserModelValidator struct {
	User struct {
		Email			string `form:"email" json:"email" binding:"exists,email"`
		FirstName		string `form:"first_name" json:"first_name" binding:"exists,min=2"`
		LastName		string `form:"last_name" json:"last_name" binding:"exists,min=2"`
		Password		string `form:"password" json:"password" binding:"exists,min=2,max=255"`
		PasswordConfirm	string `form:"password_confirm" json:"password_confirm" binding:"exists,min=2,max=255"`
	} `json:"user"`
	userModel User `json:"-"`
}

func NewUserModelValidator() UserModelValidator {
	userModelValidator := UserModelValidator{}
	return userModelValidator
}

func (u *UserModelValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, u)
	if err != nil {
		return err
	}

	u.userModel.Email = u.User.Email
	u.userModel.FirstName = u.User.FirstName
	u.userModel.LastName = u.User.LastName

	if u.User.Password != u.User.PasswordConfirm {
		return errors.New("invalid password, minimum length is 8 chars")
	}

	_ = u.userModel.SetPassword(u.User.Password)

	return nil
}