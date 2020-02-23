package authentication

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"workour-api/config"
	"workour-api/users"
)

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
	token = "the-key"
	interruptAuthentication(c, err)

	// Store session in persistence cache
	store := config.GetSessionStorage()
	profile, _ := json.Marshal(&Profile{
		Id:			user.ID,
		Email:     	user.Email,
		FirstName: 	user.FirstName,
		LastName:  	user.LastName,
	})

	// Get session with key
	session, err := store.Get(c.Request, token)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Set and save new values to session
	session.Values[token] = string(profile)
	if err = session.Save(c.Request, c.Writer); err != nil {
		log.Fatalf("Error saving session: %v", err)
	}

	// Return response with Authorization header set
	cookie := fmt.Sprintf(
		"Bearer %s; Secure; HttpOnly; SameSite=lax",
		token,
	)
	c.Writer.Header().Set("Set-cookie", cookie)
}

// Get token from authorization header
// And get user details from session stored in Redis
func (ctrl Controller) GetUserProfile(c *gin.Context) {
}

// Delete session associated with Authorization token in Redis
func (ctrl Controller) LogoutUser(c *gin.Context) {
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