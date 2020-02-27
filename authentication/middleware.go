package authentication

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"workour-api/config"
)

func VerifyAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization token from the request header first
		// and verify it is not empty
		_, err := c.Request.Cookie(CookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "No session cookie present",
			})
			return
		}

		// Get session and verify that MaxAge is not empty
		store := config.GetSessionStorage()
		session, err := store.Get(c.Request, CookieName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Error occurred while getting session",
			})
			return
		}

		if session.IsNew {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "No existing sessions were found. Please, log in.",
			})
			return
		}

		if session.Options.MaxAge == -1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Session MaxAge is -1, session ought to be terminated. Try logging in again",
			})
			return
		}
	}
}

// This handlers refreshes session by extending its expiry date
func RefreshSessionExp() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token
		// Decode it
		// Evaluate if session expires in less than ~30 minutes
		// If so, refresh session expiry, otherwise skip it
	}
}

func AuthoriseAdminAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Should verify if user accessing route is actually an admin
	}
}