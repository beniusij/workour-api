package main

import (
	"github.com/gin-gonic/gin"
	"workour-api/middleware"
	"workour-api/users"
)

func SetupRouter(driver, creds string) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.DBMiddleware(driver, creds))

	v1 := r.Group("/v1")
	User(v1.Group("/user"))

	return r
}

func User(router *gin.RouterGroup) {
	router.POST("/create", users.CreateUser)
}