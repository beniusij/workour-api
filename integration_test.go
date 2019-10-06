package main_test

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const endpoint = "/graphql"
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
	//----------------------- Test cases for getting user by id ----------------------
	{
		func(r *http.Request) {
			//resetDb(false)
		},
		endpoint,
		"POST",
		`{"query": "mutation { user: register(email: \"test@example.com\", first_name: \"Test\", last_name: \"Testest\", password: \"Password123\", password_confirm: \"Password123\") { ID } }"}`,
		http.StatusCreated,
		`{"data":{"user":{"ID":null}}}`,
		"valid data and should return StatusCreated",
	},
}

//func resetDb(addMock bool) {
//	_ = common.ResetTestDb(db)
//	db = common.InitTestDb()
//	//AutoMigrate()
//	db.AutoMigrate(&User{})
//	if addMock {
//		userMocker(10)
//	}
//}
//
//func newUserModel() User {
//	return User{
//		Email: "t3st@gmail.com",
//		FirstName: "Testas",
//		LastName: "Testavicius",
//		PasswordHash: "",
//	}
//}
//
//func userMocker(n int) []User {
//	var offset int
//	var ret []User
//	db.Model(&User{}).Count(&offset)
//
//	for i := offset + 1; i <= offset+n; i++ {
//		user := User{
//			Email: fmt.Sprintf("userModel%v@yahoo.com", i),
//			FirstName: fmt.Sprintf("User%v", i),
//			LastName: fmt.Sprintf("User%v", i),
//		}
//		_ = user.SetPassword("Password123")
//		db.Create(&user)
//		ret = append(ret, user)
//	}
//
//	return ret
//}

//func TestMain(m *testing.M) {
//	db = common.InitTestDb()
//	AutoMigrate()
//	exitval := m.Run()
//	_ = common.ResetTestDb(db)
//	os.Exit(exitval)
//}

func TestWithoutAuth(t *testing.T) {
	asserts := assert.New(t)

	r := gin.New()

	for _, testCase := range unauthRequestTestCases {
		bodyData := testCase.bodyData
		request, err := http.NewRequest(testCase.method, testCase.url, bytes.NewBufferString(bodyData))
		asserts.NoError(err)
		request.Header.Set("Content-Type", "application/json")

		testCase.init(request)

		response := httptest.NewRecorder()
		r.ServeHTTP(response, request)

		asserts.Equal(testCase.expectedCode, response.Code, "Response Status - "+testCase.msg)
		asserts.Regexp(testCase.responseRegex, response.Body.String(), "Response Content - "+testCase.msg)
	}
}
