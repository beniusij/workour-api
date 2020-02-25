package tests

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	auth "workour-api/authentication"
	"workour-api/config"
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

func TestAuthenticateUser(t *testing.T) {
	asserts := getAsserts(t)
	resetDb(true)

	// Set up router for testing
	router := gin.Default()
	config.SetupSessionStorage()

	// Login route to test login action in controller
	router.POST("/login", auth.Controller{}.AuthenticateUser)

	// Test authentication
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
				cookie := response.Header().Get("Set-Cookie")
				asserts.True(strings.Contains(cookie, "Secure"), "Cookie should be Secure")
				asserts.True(strings.Contains(cookie, "HttpOnly"), "Cookie should be HttpOnly")
				asserts.True(strings.Contains(cookie, "SameSite=Lax"), "Cookie should have SameSite=lax")
			}
		})
	}
}

func TestAuthenticatedSessionStoredInSessionStorage(t *testing.T) {
	resetDb(true)

	// Set up router for testing
	router := gin.Default()
	config.SetupSessionStorage()

	// Login route to test login action in controller
	router.POST("/login", auth.Controller{}.AuthenticateUser)

	// Stub route to test if authenticated user session is stored in session store
	router.GET("/get", func(c *gin.Context) {
		asserts := getAsserts(t)

		// Get value from store using token as key
		store := config.GetSessionStorage()
		session, err := store.Get(c.Request, auth.CookieName)

		if  err != nil {
			t.Errorf("Error found while getting session: %s", err.Error())
		}

		asserts.False(session.IsNew)
	})

	// First authenticate a user
	request, _ := http.NewRequest(
		"POST",
		"/login",
		bytes.NewBufferString(loginTestCases[0].params),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// Get token from response header
	cookie := response.Header().Get("Set-Cookie")

	// Try to get user profile from the session storage using token
	request, _ = http.NewRequest(
		"GET",
		"/get",
		nil,
	)

	request.Header.Add("Cookie", cookie)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)
}