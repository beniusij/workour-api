package authentication

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"workour-api/config"
	"workour-api/users"
)

const CookieName = "WRKSESSID"

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
}

// Authentication handler that accepts user details
func (ctrl Controller) AuthenticateUser(c *gin.Context) {
	var creds Creds
	err := json.NewDecoder(c.Request.Body).Decode(&creds)
	interruptAuthentication(c, err)

	// Get user with that email
	user, err := users.GetByEmail(creds.Email)
	interruptAuthentication(c, err)

	// Check password
	err = user.CheckPassword(creds.Password)
	interruptAuthentication(c, err)

	// Create token
	authTokenStruct := AuthToken{}
	token, err := authTokenStruct.GenerateToken(user)
	interruptAuthentication(c, err)

	// Store session in persistence cache
	store := config.GetSessionStorage()

	// Create new session with key
	session, err := store.New(c.Request, CookieName)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Set session cookie value
	session.ID = token


	// Set and save new values to session
	updateSession(session, user)
	if err = session.Save(c.Request, c.Writer); err != nil {
		log.Fatalf("Error saving session: %v", err)
	}
}

// Get token from authorization header
// And get user details from session stored in Redis
func (ctrl Controller) GetUserProfile(c *gin.Context) {
}

// Delete session associated with Authorization token in Redis
func (ctrl Controller) LogoutUser(c *gin.Context) {
}

// Update session with user details
func updateSession(s *sessions.Session, u users.User) {
	s.Values["id"] = u.ID
	s.Values["email"] = u.Email
	s.Values["first_name"] = u.FirstName
	s.Values["last_name"] = u.LastName
}

func interruptAuthentication(c *gin.Context, err error) {
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Incorrect email and/or password",
		})
		c.Abort()
		return
	}
}