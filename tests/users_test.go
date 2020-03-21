package tests

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"workour-api/config"
	"workour-api/roles"
	u "workour-api/users"
)

const regularRoleId = "Regular User"

func TestUserSettingAndCheckingPassword(t *testing.T) {
	// Set up cleaner hook
	cleaner := deleteCreatedEntities(db)
	defer cleaner()
	asserts := assert.New(t)

	// Set up test user
	user := u.User{
		Email: "t3st@gmail.com",
		FirstName: "Testas",
		LastName: "Testavicius",
		PasswordHash: "",
	}

	// Testing User password feature
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

func TestGetByEmail(t *testing.T) {
	addTestFixtures(true)

	// Set up cleaner hook
	cleaner := deleteCreatedEntities(db)
	defer cleaner()

	asserts := assert.New(t)

	email := "userModel1@yahoo.com"
	user, err := u.GetByEmail(email)

	asserts.NoError(err, "no errors are returned when fetching user by email")
	asserts.EqualValues(email, user.Email, "user fetched has the same email as the one used for getting user")

	invalidEmail := "invalid@email.com"
	user, err = u.GetByEmail(invalidEmail)

	asserts.Error(err, "record not found")
	asserts.Equal(uint(0), user.ID, "no user is returned for the invalid email")
}

func TestCreateUserResolver(t *testing.T) {
	// Set up cleaner hook
	cleaner := deleteCreatedEntities(db)
	defer cleaner()

	asserts := assert.New(t)
	userValidator := u.NewUserValidator()
	addTestFixtures(false)

	var (
		args map[string]interface{}
		err, expectedErr error
	)

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

	t.Run("creates user and returns it", func(t *testing.T) {
		args = map[string]interface{}{
			"email": 			"test1@example.com",
			"first_name":		"Test1",
			"last_name":		"Testest2",
			"password":			"Testest3!",
			"password_confirm":	"Testest3!",
		}

		err = userValidator.ValidateForm(args)
		asserts.Nil(err, "Form data validated and should not return error")

		userModel := u.User{}
		user, err := userModel.Save(userValidator.UserModel)

		// Assert response return
		asserts.Nil(err, "New userModel created with validated data")
		asserts.Equal(args["email"], user.Email, "User has email")
		asserts.Equal(getRegularUserRoleId(), user.RoleId, "Created user has default role")
	})
}

func TestGetUserResolver(t *testing.T) {
	// Set up cleaner hook
	cleaner := deleteCreatedEntities(db)
	defer cleaner()

	asserts := assert.New(t)
	userEntity := u.User{}
	addTestFixtures(true)

	t.Run("returns user with ID 1", func(t *testing.T) {
		id := uint(1)
		userEntity.ID = id

		err := userEntity.GetById()
		asserts.Nil(err, "Successfully fetched user by ID, no erros")
		asserts.Equalf(id, userEntity.ID, "Successfully fetched user with ID %v", id)
		asserts.IsType(u.User{}, userEntity, "Should return object of User interface")
	})

	t.Run("returns nil for non-existing user", func(t *testing.T) {
		id := uint(101)
		userEntity.ID = id

		err := userEntity.GetById()
		expectedErr := errors.New("record not found")
		asserts.EqualError(err, expectedErr.Error(), "Attempt to fetch non-existent user should return an error")
	})

	addTestFixtures(false)
}

// Get ID of Regular User role
func getRegularUserRoleId() uint {
	db := config.GetDB()
	role := roles.Role{Name: regularRoleId}

	err := db.First(&role).Error
	if err != nil {
		log.Println(fmt.Sprintf("Error while grabbing regular user role: %v", err))
	}

	return role.ID
}