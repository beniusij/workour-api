package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"workour-api/common"
	u "workour-api/users"
)

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
			resetDb(false)
		},
		endpoint,
		"POST",
		`{"query": "mutation { user: register(email: \"test@example.com\", first_name: \"Test\", last_name: \"Testest\", password: \"Password123\", password_confirm: \"Password123\") { ID } }"}`,
		http.StatusOK,
		`{"data":{"user":{"ID":null}}}`,
		"valid data and should return StatusCreated",
	},
}

func TestMain(m *testing.M) {
	db = common.InitTestDb()
	db.AutoMigrate(&u.User{})
	exitval := m.Run()
	_ = common.ResetTestDb(db)
	os.Exit(exitval)
}

func TestWithoutAuth(t *testing.T) {
	asserts := getAsserts(t)

	r := initTestAPI()

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
