package config

import (
	"github.com/gin-gonic/gin"
	c "workour-api/controllers"
	"workour-api/middleware"
)

func SetupRouter(driver, creds string) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.DBMiddleware(driver, creds))

	User(r.Group("/user"))


	return r
}

func User(router *gin.RouterGroup) {
	router.POST("/create", c.CreateUser)
}