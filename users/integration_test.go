package users

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"workour-api/common"
)

var db *gorm.DB
var unauthRequestTestCases = []struct{
	init			func(r *http.Request)
	url				string
	method			string
	bodyData		string
	expectedCode	int
	responseRegex	string
	msg				string
}{
	//----------------------- Test cases for userModel registration ----------------------
	{
		func(r *http.Request) {
			resetDb(false)
		},
		"/v1/user/create",
		"POST",
		`{"user":{"email":"test@example.com","first_name":"Test","last_name":"Testest","password":"TotallyValidPassword1!#","password_confirm":"TotallyValidPassword1!#"}}`,
		http.StatusCreated,
		`{"user":{"email":"test@example.com","first_name":"Test","last_name":"Testest","token":"([a-zA-Z0-9-_.]{115})"}}`,
		"valid data and should return StatusCreated",
	},
}

func resetDb(addMock bool) {
	_ = common.ResetTestDb(db)
	db = common.InitTestDb()
	//AutoMigrate()
	db.AutoMigrate(&User{})
	if addMock {
		userMocker(10)
	}
}

func newUserModel() User {
	return User{
		Email: "t3st@gmail.com",
		FirstName: "Testas",
		LastName: "Testavicius",
		PasswordHash: "",
	}
}

func userMocker(n int) []User {
	var offset int
	var ret []User
	db.Model(&User{}).Count(&offset)

	for i := offset + 1; i <= offset+n; i++ {
		user := User{
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
	AutoMigrate()
	exitval := m.Run()
	_ = common.ResetTestDb(db)
	os.Exit(exitval)
}

func TestUser(t *testing.T) {
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

func TestWithoutAuth(t *testing.T) {
	asserts := assert.New(t)

	r := gin.New()
	UserRoutes(r.Group("/v1"))

	for _, testCase := range unauthRequestTestCases {
		bodyData := testCase.bodyData
		req, err := http.NewRequest(testCase.method, testCase.url, bytes.NewBufferString(bodyData))
		asserts.NoError(err)
		req.Header.Set("Content-Type", "application/json")

		testCase.init(req)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		asserts.Equal(testCase.expectedCode, w.Code, "Response Status - "+testCase.msg)
		asserts.Regexp(testCase.responseRegex, w.Body.String(), "Response Content - "+testCase.msg)
	}
}
