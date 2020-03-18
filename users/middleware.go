package users

//import (
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"log"
//	"workour-api/authentication"
//	"workour-api/config"
//)
//
//func LoadUser(c *gin.Context)  {
//	// Get cookie
//	store := config.GetSessionStorage()
//	session, err := store.Get(c.Request, authentication.CookieName)
//	if err != nil {
//		log.Println(fmt.Sprintf("Failed to load current session: %v", err))
//		c.Next()
//	}
//
//	// Load user
//	user := User{ID: session.Values["id"].(uint)}
//	err = user.GetById()
//	if err != nil {
//		log.Println(fmt.Sprintf("Failed to load user: %v", err))
//		c.Next()
//	}
//
//	// Store in context
//	c.Keys["user"] = user
//}