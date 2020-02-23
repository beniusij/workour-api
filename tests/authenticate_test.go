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

func TestAuthenticateUser(t *testing.T) {
	asserts := getAsserts(t)
	resetDb(true)

	// Set up router for testing
	router := gin.Default()
	common.InitSessionStore(router)

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
				cookie := response.Header().Get("Set-cookie")
				asserts.True(strings.Contains(cookie, tokenType))
				asserts.True(strings.Contains(cookie, "Secure"))
				asserts.True(strings.Contains(cookie, "HttpOnly"))
				asserts.True(strings.Contains(cookie, "SameSite=lax"))
			}
		})
	}
}

//func TestAuthenticatedSessionStoredInSessionStorage(t *testing.T) {
//	//asserts := getAsserts(t)
//	resetDb(true)
//
//	// Set up router for testing
//	router := gin.Default()
//	common.InitSessionStore(router)
//
//	// Login route to test login action in controller
//	router.POST("/login", auth.Controller{}.AuthenticateUser)
//
//	// Stub route to test if authenticated user session is stored in session store
//	router.GET("/get", func(c *gin.Context) {
//		// Get token from request body
//		token := c.Request.URL.Query().Get("token")
//
//		// Get value from store using token as key
//		session := sessions.Default(c)
//		fmt.Println(token)
//		p := session.Get(token)
//		if  p == nil {
//			t.Errorf("Not found for %s", token)
//		}
//		session.Save()
//	})
//
//	// First authenticate a user
//	request, _ := http.NewRequest(
//		"POST",
//		"/login",
//		bytes.NewBufferString(loginTestCases[0].params),
//	)
//	request.Header.Set("Content-Type", "application/json")
//	response := httptest.NewRecorder()
//	router.ServeHTTP(response, request)
//
//	// Get token from response header
//	cookie := response.Header().Get("Set-cookie")
//	splitCookie := strings.Split(cookie, " ")
//	token := strings.Trim(splitCookie[1], ";")
//
//	// Try to get user profile from the session storage using token
//	request, _ = http.NewRequest(
//		"GET",
//		fmt.Sprintf("/get?token=%s", token),
//		nil,
//	)
//	response = httptest.NewRecorder()
//	router.ServeHTTP(response, request)
//}