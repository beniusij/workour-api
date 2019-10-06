package users_test

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"workour-api/common"
	u "workour-api/users"
)

var db *gorm.DB

func resetDb(addMock bool) {
	_ = common.ResetTestDb(db)
	db = common.InitTestDb()
	//AutoMigrate()
	db.AutoMigrate(&u.User{})
	if addMock {
		userMocker(10)
	}
}

func newUserModel() u.User {
	return u.User{
		Email: "t3st@gmail.com",
		FirstName: "Testas",
		LastName: "Testavicius",
		PasswordHash: "",
	}
}

func userMocker(n int) []u.User {
	var offset int
	var ret []u.User
	db.Model(&u.User{}).Count(&offset)

	for i := offset + 1; i <= offset+n; i++ {
		user := u.User{
			Email: fmt.Sprintf("userModel%v@yahoo.com", i),
			FirstName: fmt.Sprintf("User%v", i),
			LastName: fmt.Sprintf("User%v", i),
		}
		_ = user.SetPassword("Password123")
		db.Create(&user)
		ret = append(ret, user)
	}

	return ret
}

func TestMain(m *testing.M) {
	db = common.InitTestDb()
	db.AutoMigrate(&u.User{})
	exitval := m.Run()
	_ = common.ResetTestDb(db)
	os.Exit(exitval)
}

func TestUserSettingAndCheckingPassword(t *testing.T) {
	asserts := assert.New(t)

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
	asserts := assert.New(t)
	userValidator := u.NewUserValidator()

	args := map[string]interface{}{
		"email": 			"test1",
		"first_name":		"Te",
		"last_name":		"",
		"password":			"Te",
		"password_confirm":	"Test",
	}

	err := userValidator.ValidateForm(args)
	expectedErr := errors.New("an error is expected but got nil")
	asserts.EqualError(err, expectedErr.Error(), "Form data did not validate and returns an error")

	args = map[string]interface{}{
		"email": 			"test1@example.com",
		"first_name":		"Test1",
		"last_name":		"Testest2",
		"password":			"Testest3!",
		"password_confirm":	"Testest3!",
	}
	t.Logf("%v", args)

	err = userValidator.ValidateForm(args)
	asserts.Nil(err, "Form data validated and should not return error")

	err = u.SaveUser(&userValidator.UserModel)
	asserts.Nil(err, "New user created with validated data")
	fmt.Println(fmt.Sprintf("%v", &userValidator.UserModel.ID))
	asserts.NotNil(userValidator.UserModel.ID, "User has unique ID")
}