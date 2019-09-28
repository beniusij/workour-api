package config

import (
	"github.com/gin-gonic/gin"
	u "workour-api/users"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	u.UserRoutes(v1.Group("/user"))

	return r
}

