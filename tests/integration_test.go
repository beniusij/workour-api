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
	//----------------------- Test cases for registering user ----------------------
	{
		func(r *http.Request) {
			resetDb(false)
		},
		endpoint,
		"POST",
		`{"query": "mutation { user: register(email: \"test@example.com\", first_name: \"Test\", last_name: \"Testest\", password: \"Password123\", password_confirm: \"Password123\") { ID } }"}`,
		http.StatusCreated,
		`{"data":{"user":{"ID":1}}}`,
		"valid data and should return StatusCreated",
	},
	{
		func(r *http.Request) {},
		endpoint,
		"POST",
		`{"query": "mutation { user: register(email: \"test@example.com\", first_name: \"Test\", last_name: \"Testest\", password: \"Password123\", password_confirm: \"Password123\") { ID } }"}`,
		http.StatusBadRequest,
		"UNIQUE constraint failed: users.email",
		"use of non-unique email should fail and return StatusBadRequest",
	},
	{
		func(r *http.Request) {},
		endpoint,
		"POST",
		`{"query": "mutation { user: register(email: \"test1\", first_name: \"Test\", last_name: \"Testest\", password: \"Password123\", password_confirm: \"Password123\") { ID } }"}`,
		http.StatusBadRequest,
		"Error:Field validation for 'Email' failed on the 'email' tag",
		"form with invalid email should fail and return StatusBadRequest",
	},
	{
		func(r *http.Request) {},
		endpoint,
		"POST",
		`{"query": "mutation { user: register(email: \"test1@example.com\", first_name: \"T\", last_name: \"\", password: \"Password123\", password_confirm: \"\") { ID } }"}`,
		http.StatusBadRequest,
		"min",
		"form with invalid first & last names, password and confirm password fields should fail and return StatusBadRequest",
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

	for _, tc := range unauthRequestTestCases {
		t.Run(tc.msg, func(t *testing.T) {
			bodyData := tc.bodyData
			request, err := http.NewRequest(tc.method, tc.url, bytes.NewBufferString(bodyData))
			asserts.NoError(err)
			request.Header.Set("Content-Type", "application/json")

			tc.init(request)

			response := httptest.NewRecorder()
			r.ServeHTTP(response, request)

			asserts.Equal(tc.expectedCode, response.Code, "Response Status - "+tc.msg)
			asserts.Regexp(tc.responseRegex, response.Body.String(), "Response Content - "+tc.msg)
		})
	}
}
