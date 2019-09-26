package main

import (
	"github.com/gin-gonic/gin"
	"workour-api/users"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	User(v1.Group("/user"))

	return r
}

func User(router *gin.RouterGroup) {
	router.POST("/create", users.CreateUser)
}