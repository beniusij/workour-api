package validators

import (
	"errors"
	"github.com/gin-gonic/gin"
	"workour-api/helpers"
	"workour-api/models"
)

type UserModelValidator struct {
	User struct {
		Email			string `form:"email" json:"email" binding:"exists,email"`
		FirstName		string `form:"first_name" json:"first_name" binding:"exists,min=2"`
		LastName		string `form:"last_name" json:"last_name" binding:"exists,min=2"`
		Password		string `form:"password" json:"password" binding:"exists,min=2,max=255"`
		PasswordConfirm	string `form:"password" json:"password" binding:"exists,min=2,max=255"`
	} `json:"user"`
	user models.User `json:"-"`
}

func NewUserModelValidator() UserModelValidator {
	return UserModelValidator{}
}

func (u *UserModelValidator) Bind(c *gin.Context) error {
	err := helpers.Bind(c, u)
	if err != nil {
		return err
	}

	u.user.Email = u.User.Email
	u.user.FirstName = u.User.FirstName
	u.user.LastName = u.User.LastName

	if u.User.Password != u.User.PasswordConfirm {
		return errors.New("invalid password, minimum length is 8 chars")
	}

	_ = u.user.SetPassword(u.User.Password)

	return nil
}