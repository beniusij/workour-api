package tests

import (
	"errors"
	"testing"
	u "workour-api/users"
)

func TestUserSettingAndCheckingPassword(t *testing.T) {
	asserts := getAsserts(t)

	// Testing User password feature
	user := newUserModel()
	err := user.CheckPassword("")
	asserts.Error(err, "empty password should return err")

	err = user.SetPassword("Password12355!")
	asserts.NoError(err, "password should be set successful")
	asserts.Len(user.PasswordHash, 60, "password hash length should be 60")

	err = user.CheckPassword("Password12355")
	asserts.Error(err, "password should be checked and not validated")

	err = user.CheckPassword("Password12355!")
	asserts.NoError(err, "password should be checked and validated")
}

func TestCreateUserResolver(t *testing.T) {
	asserts := getAsserts(t)
	userValidator := u.NewUserValidator()
	var (
		args map[string]interface{}
		err, expectedErr error
	)
	resetDb(false)

	// Init faulty testing data
	var faultyData = []struct{
		msg		string
		args	map[string]interface{}
	}{
		{
			"standard faulty data, returns JSON with errors",
			map[string]interface{}{
				"email": 			"test1",
				"first_name":		"Te",
				"last_name":		"",
				"password":			"Te",
				"password_confirm":	"Test",
			},
		},
		{
			"only numbers in data, returns JSON with errors",
			map[string]interface{}{
				"email": 			"123456@12314.12",
				"first_name":		"12",
				"last_name":		"",
				"password":			"34",
				"password_confirm":	"456678990",
			},
		},
		{
			"data with html tags and symbols, returns JSON with errors",
			map[string]interface{}{
				"email": 			"<input type='email'>test@example.com</input>",
				"first_name":		"<p>My Test Name</p>",
				"last_name":		"",
				"password":			"<>",
				"password_confirm":	"~!@£$%^&*()_+|}{P:?><",
			},
		},
	}

	for _, data := range faultyData {
		t.Run(data.msg, func(t *testing.T) {
			err = userValidator.ValidateForm(data.args)
			expectedErr = errors.New("Key: 'Email' Error:Field validation for 'Email' failed on the 'email' tag\nKey: 'LastName' Error:Field validation for 'LastName' failed on the 'required' tag\nKey: 'Password' Error:Field validation for 'Password' failed on the 'min' tag\nKey: 'PasswordConfirm' Error:Field validation for 'PasswordConfirm' failed on the 'eqfield' tag")
			asserts.EqualError(err, expectedErr.Error(), "Form data did not validate and returns an error")
		})
	}

	t.Run("creates user and returns its ID", func(t *testing.T) {
		args = map[string]interface{}{
			"email": 			"test1@example.com",
			"first_name":		"Test1",
			"last_name":		"Testest2",
			"password":			"Testest3!",
			"password_confirm":	"Testest3!",
		}

		err = userValidator.ValidateForm(args)
		asserts.Nil(err, "Form data validated and should not return error")

		user := u.User{}
		var id int
		id, err = user.SaveEntity(&userValidator.UserModel)
		asserts.Nil(err, "New user created with validated data")
		asserts.Equal(1, id, "User has ID 1")
	})

	resetDb(false)
}

func TestGetUserResolver(t *testing.T) {
	asserts := getAsserts(t)
	userEntity := u.User{}
	userMocker(10)
	var (
		id		int
		args 	map[string]interface{}
		err		error
		user	*u.User
	)

	t.Run("returns user with ID 1", func(t *testing.T) {
		id = 1
		args = map[string]interface{}{
			"id":	id,
		}

		user, err = userEntity.GetEntityById(args["id"].(int))
		asserts.Nil(err, "Successfully fetched user by ID, no erros")
		asserts.Equalf(id, user.ID, "Successfully fetched user with ID %v", id)
	})

	t.Run("returns nil for non-existing user", func(t *testing.T) {
		id = 101
		args = map[string]interface{}{
			"id": id,
		}

		user, err = userEntity.GetEntityById(args["id"].(int))
		expectedErr := errors.New("record not found")
		asserts.Nil(user, "Attempt to fetch non-existent user returns nil for user")
		asserts.EqualError(err, expectedErr.Error(), "Attempt to fetch non-existent user should return an error")
	})

	resetDb(false)
}