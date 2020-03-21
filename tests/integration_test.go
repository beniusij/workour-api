package tests

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"workour-api/config"
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
			_ = config.ResetTestDb(db)
			db = config.InitTestDb()
			migrate()

			roleMocker()
		},
		publicEndpoint,
		"POST",
		`{"query": "mutation { user: register(email: \"test@example.com\", first_name: \"Test\", last_name: \"Testest\", password: \"Password123\", password_confirm: \"Password123\") { ID } }"}`,
		http.StatusCreated,
		`{"data":{"user":{"ID":1}}}`,
		"valid data and should return StatusCreated",
	},
	{
		func(r *http.Request) {},
		publicEndpoint,
		"POST",
		`{"query": "mutation { user: register(email: \"test@example.com\", first_name: \"Test\", last_name: \"Testest\", password: \"Password123\", password_confirm: \"Password123\") { ID } }"}`,
		http.StatusOK,
		"UNIQUE constraint failed: users.email",
		"use of non-unique email should fail and return StatusBadRequest",
	},
	{
		func(r *http.Request) {},
		publicEndpoint,
		"POST",
		`{"query": "mutation { user: register(email: \"test1\", first_name: \"Test\", last_name: \"Testest\", password: \"Password123\", password_confirm: \"Password123\") { ID } }"}`,
		http.StatusOK,
		"Error:Field validation for 'Email' failed on the 'email' tag",
		"form with invalid email should fail and return StatusBadRequest",
	},
	{
		func(r *http.Request) {},
		publicEndpoint,
		"POST",
		`{"query": "mutation { user: register(email: \"test1@example.com\", first_name: \"T\", last_name: \"Testst\", password: \"Password123\", password_confirm: \"Password123\") { ID } }"}`,
		http.StatusOK,
		`Field validation for 'FirstName' failed on the 'min' tag"`,
		"form with invalid first name should fail and return StatusBadRequest",
	},
	{
		func(r *http.Request) {},
		publicEndpoint,
		"POST",
		`{"query": "mutation { user: register(email: \"test1@example.com\", first_name: \"Test\", last_name: \"\", password: \"Password123\", password_confirm: \"Password123\") { ID } }"}`,
		http.StatusOK,
		`Field validation for 'LastName' failed on the 'required' tag"`,
		"form with no last name should fail and return StatusBadRequest",
	},
	{
		func(r *http.Request) {},
		publicEndpoint,
		"POST",
		`{"query": "mutation { user: register(email: \"test1@example.com\", first_name: \"Test\", last_name: \"Testest\", password: \"Password123\", password_confirm: \"Password12\") { ID } }"}`,
		http.StatusOK,
		`Field validation for 'PasswordConfirm' failed on the 'eqfield' tag"`,
		"form with not matching passwords should fail and return StatusBadRequest",
	},
}

func TestMain(m *testing.M) {
	// Pre-testing setup
	db = config.InitTestDb()
	migrate()

	exitValue := m.Run()

	// Post-testing cleanup
	_ = config.ResetTestDb(db)

	os.Exit(exitValue)
}

func TestWithoutAuth(t *testing.T) {
	// Set up cleaner hook
	cleaner := deleteCreatedEntities(db)
	defer cleaner()
	
	asserts := assert.New(t)
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
