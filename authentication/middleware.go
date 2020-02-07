package authentication

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func VerifyAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization token from the request header first
		// and verify it is not empty
		authToken := c.Request.Header.Get("Authorization")
		if authToken != "" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "No authentication token present",
			})
			c.Abort()
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