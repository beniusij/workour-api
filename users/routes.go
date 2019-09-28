package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRoutes(router *gin.RouterGroup) {
	r := router.Group("/user")

	r.POST("/create", CreateUser)
	r.GET("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, "{'action': 'read'}")
	})
	r.PUT("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, "{'action': 'update}")
	})
	r.DELETE("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, "{'action': 'delete'}")
	})
}