package authentication

import "github.com/gin-gonic/gin"

type Controller struct {}

// Authentication handler that accepts user details
func (ctrl Controller) AuthenticateUser(c *gin.Context) {
		// Get user with that email
		// Check password
		// Create token
		// Return response with Authorization header set
}

// Get token from authorization header
// And get user details from session stored in Redis
func (ctrl Controller) GetUserProfile(c *gin.Context) {
}

// Delete session associated with Authorization token
func (ctrl Controller) LogoutUser(c *gin.Context) {
}
