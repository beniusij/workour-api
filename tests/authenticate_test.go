package tests

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	auth "workour-api/authentication"
	"workour-api/common"
)


const tokenType = "Bearer"
var loginTestCases = []struct{
	msg		string
	params	string
	status	int
}{
	{
		"Should authenticate user",
		`{"email":"userModel1@yahoo.com","password":"Password123"}`,
		200,
	},
	{
		"Should not authenticate user with incorrect email",
		`{"email":"visata@gmail.com","password":"Password123"}`,
		200,
	},
	{
		"Should not authenticate user with incorrect password",
		`{"email":"userModel1@yahoo.com","password":"Informatika"}`,
		200,
	},
}

func TestAuthenticateUSer(t *testing.T) {
	asserts := getAsserts(t)
	resetDb(true)

	// Set up router for testing
	router := gin.Default()
	common.InitSessionStore(router, "localhost:6379")
	router.POST("/login", auth.Controller{}.AuthenticateUser)

	for _, testCase := range loginTestCases {
		t.Run(testCase.msg, func(t *testing.T) {
			request, err := http.NewRequest(
				"POST",
				"/login",
				bytes.NewBufferString(testCase.params),
			)
			asserts.NoError(err)
			request.Header.Set("Content-Type", "application/json")

			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			asserts.Equal(testCase.status, response.Code, "Correct status code is received")

			body := response.Body.String()
			if !strings.Contains(body, "Incorrect email and/or password") {
				cookie := response.Header().Get("Set-cookie")
				asserts.True(strings.Contains(cookie, tokenType))
				asserts.True(strings.Contains(cookie, "Secure"))
				asserts.True(strings.Contains(cookie, "HttpOnly"))
				asserts.True(strings.Contains(cookie, "SameSite=lax"))
			}
		})
	}
}