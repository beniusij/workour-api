package users

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
)

type UserModelValidator struct {
	user struct{
		Email			string `validate:"required,email"`
		FirstName		string `validate:"required,min=2,max=255"`
		LastName		string `validate:"required,min=2,max=255"`
		Password		string `validate:"required,min=8,max=255"`
		PasswordConfirm	string `validate:"required,eqfield=Password"`
	}
	UserModel User
}

var validate *validator.Validate

func NewUserValidator() UserModelValidator {
	userValidator := UserModelValidator{}
	return userValidator
}

func (u *UserModelValidator) ValidateForm(p map[string]interface{}) error {
	validate = validator.New()

	// Unmarshal params from graphql.ResolveParams and put in struct for validation
	u.user.Email = p["email"].(string)
	u.user.FirstName = p["first_name"].(string)
	u.user.LastName = p["last_name"].(string)
	u.user.Password = p["password"].(string)
	u.user.PasswordConfirm = p["password_confirm"].(string)

	err := validate.Struct(u.user)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error validating form: %v", err))
		return err
	}

	// After validation re-assign those values to User and set password hash
	u.UserModel.Email = u.user.Email
	u.UserModel.FirstName = u.user.FirstName
	u.UserModel.LastName = u.user.LastName
	err = u.UserModel.SetPassword(u.user.Password)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error setting password: %v", err))
		return err
	}

	return nil
}