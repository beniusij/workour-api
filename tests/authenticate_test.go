package tests

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	auth "workour-api/authentication"
	"workour-api/config"
)

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
		403,
	},
	{
		"Should not authenticate user with incorrect password",
		`{"email":"userModel1@yahoo.com","password":"Informatika"}`,
		403,
	},
}

func TestAuthenticateUser(t *testing.T) {
	asserts, router := setUp(t)

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
				asserts.True(strings.Contains(cookie, "HttpOnly"), "Cookie should be HttpOnly")
				asserts.True(strings.Contains(cookie, "SameSite=Lax"), "Cookie should have SameSite=lax")
			}
		})
	}
}

func TestAuthenticatedSessionStoredInSessionStorage(t *testing.T) {
	_, router := setUp(t)

	// Add login route with handler and authenticate as a test user
	cookie := authTestUser(router)

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
		asserts.Equal("userModel1@yahoo.com", session.Values["email"])
	})

	// Try to get user profile from the session storage using token
	request, _ := http.NewRequest(
		"GET",
		"/get",
		nil,
	)

	request.Header.Add("Cookie", cookie)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
}

func TestUserCanLogout(t *testing.T) {
	asserts, router := setUp(t)
	store := config.GetSessionStorage()

	cookie := authTestUser(router)

	// Add routes
	router.POST("/logout", auth.Controller{}.LogoutUser)

	// Prep for logout
	request, _ := http.NewRequest("POST", "/logout", nil)
	request.Header.Add("Cookie", cookie)
	response := httptest.NewRecorder()

	// Check that session is created
	session, err := store.Get(request, auth.CookieName)
	logErr(err)

	asserts.False(session.IsNew, "Session is created and it is in Redis")
	asserts.Equal("userModel1@yahoo.com", session.Values["email"])

	// Log out user
	router.ServeHTTP(response, request)

	// Check that session is created
	session, err = store.Get(request, auth.CookieName)
	logErr(err)

	asserts.Equal(
		-1,
		session.Options.MaxAge,
		"Session is destroyed and user is no longer authenticated",
		)
}

func TestSessionAuthenticationMiddleware(t *testing.T) {
	asserts, router := setUp(t)

	router.POST("/logout", auth.Controller{}.LogoutUser)
	router.GET("/ping", auth.VerifyAuthentication, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Request received",
		})
	})

	// TODO: refactor into slice of test cases
	//testCase := []struct{
	//	case string
	//}{}

	// Make a call to the protected route without a cookie in the request
	request, err := http.NewRequest("GET", "/ping", nil)
	logErr(err)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// Should fail and return response 403 with error "No authentication cookie
	// present"
	asserts.Equal(http.StatusForbidden, response.Code, "No valid cookie, thus access is forbidden")
	asserts.True(strings.Contains(response.Body.String(), "No session cookie present"))

	// Send request with fake cookie
	request.Header.Add("Cookie", fmt.Sprintf("%s=FakeCookie", auth.CookieName))
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// Should return response 403 with error "No existing sessions were found.
	// Please, log in."
	asserts.Equal(http.StatusForbidden, response.Code, "Invalid cookie, thus access is forbidden")
	asserts.True(strings.Contains(response.Body.String(), "Error occurred while getting session"))

	// Authenticate as test user
	cookie := authTestUser(router)

	// Log out
	request, _ = http.NewRequest("POST", "/logout", nil)
	request.Header.Add("Cookie", cookie)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// Send request to the protected route
	request, _ = http.NewRequest("GET", "/ping", nil)
	request.Header.Add("Cookie", cookie)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// Should return response 403 with error "Session MaxAge is -1, session ought
	// to be terminated. Try logging in again"
	asserts.Equal(http.StatusForbidden, response.Code, "Invalid cookie, thus access is forbidden")
	asserts.True(
		strings.Contains(
			response.Body.String(),
			"No existing sessions were found. Please, log in.",
		),
	)

	// Authenticate as test user
	authTestUser(router)

	// Send request to protected route
	request, _ = http.NewRequest("GET", "/ping", nil)
	request.Header.Add("Cookie", cookie)
	response = httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// Should return response with HTTP code 200
	asserts.Equal(http.StatusOK, response.Code, "Valid cookie, thus access is not prevented")
	asserts.True(strings.Contains(response.Body.String(), "Request received"))
}

func TestLoadUserMiddleware(t *testing.T) {
	asserts, router := setUp(t)

	router.GET("/ping", auth.LoadUser, func(c *gin.Context) {
		user := c.Keys["user"]

		if user != nil {
			c.JSON(http.StatusOK, gin.H{
				"currentUser": user,
			})
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Request received",
			})
		}
	})

	cookie := authTestUser(router)

	request, _ := http.NewRequest("GET", "/ping", nil)
	request.Header.Add("Cookie", cookie)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// Assert response body
	asserts.Equal(http.StatusOK, response.Code, "Should return response 200")
}

func TestGetCurrentUser(t *testing.T) {
	asserts, router := setUp(t)

	// Log in as user
	cookie := authTestUser(router)

	router.GET("/getCurrentUser", auth.Controller{}.GetCurrentUser)

	// Get current authenticated user
	request, _ := http.NewRequest("GET", "/getCurrentUser", nil)
	request.Header.Add("Cookie", cookie)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	// Assert result
	result := response.Body.String()
	asserts.Equal(http.StatusOK, response.Code, "Response with HTTP code 200 returned")
	asserts.True(strings.Contains(result, "userModel1@yahoo.com"))

	// Assert it has correct role id
	roleId := getRegularUserRoleId()
	expected := fmt.Sprintf(`\"Role\":{\"ID\":%d`, roleId)
	asserts.Contains(result, expected, "Current user json has correct role id")
}

func logErr(err error) {
	if err != nil {
		log.Println(fmt.Sprintf("Error occurred while running test: %v", err))
	}
}

// Adds login route if missing and authenticates as test user
func authTestUser(r *gin.Engine) string {
	if !hasRoute(r, "/login") {
		r.POST("/login", auth.Controller{}.AuthenticateUser)
	}

	testUser := `{"email":"userModel1@yahoo.com","password":"Password123"}`

	request, _ := http.NewRequest(
		"POST",
		"/login",
		bytes.NewBufferString(testUser),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	r.ServeHTTP(response, request)

	return response.Header().Get("Set-Cookie")
}

// Helper method to check if router has route registered
func hasRoute(r *gin.Engine, path string) bool {
	routes := r.Routes()

	for _, route := range routes {
		if route.Path == path {
			return true
		}
	}

	return false
}

func setUp(t *testing.T) (a *assert.Assertions, r *gin.Engine) {
	resetDb(true)
	config.SetupSessionStorage()

	a = getAsserts(t)
	r = gin.New()

	return
}