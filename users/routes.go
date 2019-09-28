package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserRoutes(router *gin.RouterGroup) {
	router.POST("/create", CreateUser)
	router.GET("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, "{'action': 'read'}")
	})
	router.PUT("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, "{'action': 'update}")
	})
	router.DELETE("/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, "{'action': 'delete'}")
	})
}