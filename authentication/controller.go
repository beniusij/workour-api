package authentication

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
	"workour-api/config"
	"workour-api/roles"
	"workour-api/users"
)

const CookieName = "WRKSESSID"

var env = os.Getenv("APP_ENV")
var secureFlag = true
var domain = os.Getenv("DOMAIN")

type Controller struct {}
type Creds 		struct {
	Email 		string
	Password 	string
}
type Profile 	struct {
	Id			uint
	Email 		string
	FirstName 	string
	LastName 	string
	Role		roles.Role
}

// Authentication handler that accepts user details
func (ctrl Controller) AuthenticateUser(c *gin.Context) {
	var creds Creds
	err := json.NewDecoder(c.Request.Body).Decode(&creds)
	if err != nil {
		interruptAuthentication(c, err)
		return
	}

	// Get user with that email
	user, err := users.GetByEmail(creds.Email)
	if err != nil {
		interruptAuthentication(c, err)
		return
	}

	// Check password
	err = user.CheckPassword(creds.Password)
	if err != nil {
		interruptAuthentication(c, err)
		return
	}

	// Create token
	authTokenStruct := AuthToken{}
	token, err := authTokenStruct.GenerateToken(user)
	if err != nil {
		interruptAuthentication(c, err)
		return
	}

	// Store session in persistence cache
	store := config.GetSessionStorage()

	// Create new session with key
	session, err := store.New(c.Request, CookieName)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Set session cookie value
	session.ID = token
	session.Options = setCookieOptions()
	
	// Set and save new values to session
	updateSession(session, user)
	if err = session.Save(c.Request, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not save session",
		})
		c.Abort()
		return
	}
	c.Status(http.StatusOK)
}

// Get user details from session stored in Redis
func (ctrl Controller) GetCurrentUser(c *gin.Context) {
	store := config.GetSessionStorage()
	session, err := store.Get(c.Request, CookieName)
	if err != nil {
		log.Println(fmt.Sprintf("Error occurred while getting current user: %v", err))
	}

	roleId := session.Values["role_id"].(uint)

	profile := Profile{
		Id: session.Values["id"].(uint),
		Email: session.Values["email"].(string),
		FirstName: session.Values["first_name"].(string),
		LastName: session.Values["last_name"].(string),
		Role: roles.GetRoleById(roleId),
	}

	profileJson, err := json.Marshal(&profile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to parse profile into JSON",
		})
	}

	c.JSON(http.StatusOK, string(profileJson))
}

// Delete session associated with Authorization token in Redis
func (ctrl Controller) LogoutUser(c *gin.Context) {
	store := config.GetSessionStorage()

	// Get session
	session, err := store.Get(c.Request, CookieName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Could not get session",
		})
		c.Abort()
	}

	// Set session.Options.MaxAge = -1 and save to delete the session
	session.Options.MaxAge = -1
	if err = session.Save(c.Request, c.Writer); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Could not save session",
		})
		c.Abort()
		return
	}
}

// Update session with user details
func updateSession(s *sessions.Session, u users.User) {
	s.Values["id"] = u.ID
	s.Values["email"] = u.Email
	s.Values["first_name"] = u.FirstName
	s.Values["last_name"] = u.LastName
	s.Values["role_id"] = u.RoleId
}

func interruptAuthentication(c *gin.Context, err error) {
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Incorrect email and/or password",
		})
		c.Abort()
	}
}

// Sets options for cookie which is later passed to session
func setCookieOptions() *sessions.Options {
	if env == "development" {
		secureFlag = false
		domain = "localhost"
	}

	return &sessions.Options{
		Path:     "/",
		Domain:   domain,
		MaxAge:   3600 * 24,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}