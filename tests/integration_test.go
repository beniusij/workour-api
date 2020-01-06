package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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
		http.StatusBadRequest,
		"UNIQUE constraint failed: users.email",
		"use of non-unique email should fail and return StatusBadRequest",
	},
	{
		func(r *http.Request) {},
		publicEndpoint,
		"POST",
		`{"query": "mutation { user: register(email: \"test1\", first_name: \"Test\", last_name: \"Testest\", password: \"Password123\", password_confirm: \"Password123\") { ID } }"}`,
		http.StatusBadRequest,
		"Error:Field validation for 'Email' failed on the 'email' tag",
		"form with invalid email should fail and return StatusBadRequest",
	},
	{
		func(r *http.Request) {},
		publicEndpoint,
		"POST",
		`{"query": "mutation { user: register(email: \"test1@example.com\", first_name: \"T\", last_name: \"Testst\", password: \"Password123\", password_confirm: \"Password123\") { ID } }"}`,
		http.StatusBadRequest,
		`Field validation for 'FirstName' failed on the 'min' tag"`,
		"form with invalid first name should fail and return StatusBadRequest",
	},
	{
		func(r *http.Request) {},
		publicEndpoint,
		"POST",
		`{"query": "mutation { user: register(email: \"test1@example.com\", first_name: \"Test\", last_name: \"\", password: \"Password123\", password_confirm: \"Password123\") { ID } }"}`,
		http.StatusBadRequest,
		`Field validation for 'LastName' failed on the 'required' tag"`,
		"form with no last name should fail and return StatusBadRequest",
	},
	{
		func(r *http.Request) {},
		publicEndpoint,
		"POST",
		`{"query": "mutation { user: register(email: \"test1@example.com\", first_name: \"Test\", last_name: \"Testest\", password: \"Password123\", password_confirm: \"Password12\") { ID } }"}`,
		http.StatusBadRequest,
		`Field validation for 'PasswordConfirm' failed on the 'eqfield' tag"`,
		"form with not matching passwords should fail and return StatusBadRequest",
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

const tokenType = "Bearer"
var loginTestCases = []struct{
	url				string
	msg				string
	query			string
	expectedCode	int
	Email			string
	isAuth			bool
}{
	{
		publicEndpoint,
		"form with correct credentials should return response with JWT token and user's name and email",
		`{"query": "mutation { user: login(email: \"userModel1@yahoo.com\", password: \"Password123\") { Email, Token } }"}`,
		200,
		"userModel1@yahoo.com",
		true,
	},
	{
		publicEndpoint,
		"form with incorrect credentials should return response without user and JWT",
		`{"query": "mutation { user: login(email: \"userModel2@yahoo.com\", password: \"<p disabled>Incorrect psw</p>\") { Email, Token } }"}`,
		200,
		"",
		false,
	},
	{
		publicEndpoint,
		"form with incorrect credentials should return response without user and JWT",
		`{"query": "mutation { user: login(email: \"\", password: \"\") { Email, Token } }"}`,
		200,
		"",
		false,
	},
}


// Struct to unmarshal response body
type RespData struct {
	Data struct {
		User struct {
			Email string
			Token string
		} `json:"user"`
	} `json:"data"`
}

func TestAuthentication(t *testing.T) {
	asserts := getAsserts(t)
	r := initTestAPI()
	resetDb(true)
	var resp RespData

	for _, tc := range loginTestCases {
		t.Run(tc.msg, func(t *testing.T) {
			query := tc.query
			request, err := http.NewRequest("POST", tc.url, bytes.NewBufferString(query))
			asserts.NoError(err)
			request.Header.Set("Content-Type", "application/json")

			response := httptest.NewRecorder()
			r.ServeHTTP(response, request)

			// Assert response status code
			asserts.Equal(tc.expectedCode, response.Code, fmt.Sprintf("Response Status - %s", tc.msg))

			respBody := response.Body.String()
			if strings.Contains(respBody, `"data":{"user":null}`) {
				asserts.Contains(respBody, "incorrect email and/or password", "Contains correct error message")

				resp = RespData{}
			} else {
				err = json.Unmarshal(response.Body.Bytes(), &resp)
				asserts.NoError(err)

				asserts.Equal(tc.Email, resp.Data.User.Email, fmt.Sprintf("Response Content - %s", tc.msg))
				asserts.Regexp(TokenRegex, resp.Data.User.Token, "JWT token matches token regex")
			}
		})
	}
}
