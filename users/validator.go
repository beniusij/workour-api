package users

import (
	"errors"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
)

type UserValidator struct { 
	Email			string `validate:"required,email,unique"`
	FirstName		string `validate:"required,min=2,max=255"`
	LastName		string `validate:"required,min=2,max=255"`
	Password		string `validate:"required,min=8,max=255"`
	PasswordConfirm	string `validate:"required,eqfield=Password"`
}

var validate *validator.Validate

func validateUserForm(p map[string]interface{}) (interface{}, error) {
	validate = validator.New()

	user := &UserValidator{
		Email: p["email"].(string),
		FirstName: p["first_name"].(string),
		LastName: p["last_name"].(string),
		Password: p["password"].(string),
		PasswordConfirm: p["password_confirm"].(string),
	}

	if err := validate.Struct(user); err != nil {
		err := errors.New(fmt.Sprintf("submitted user registration data did not pass validation: %v", err))
		return nil, err
	}

	//err := common.Bind(c, u)
	//if err != nil {
	//	return err
	//}
	//
	//u.userModel.Email = u.User.Email
	//u.userModel.FirstName = u.User.FirstName
	//u.userModel.LastName = u.User.LastName
	//
	//if u.User.Password != u.User.PasswordConfirm {
	//	return errors.New("invalid password, minimum length is 8 chars")
	//}
	//
	//_ = u.userModel.SetPassword(u.User.Password)

	return user, nil
}