package main

import (
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
	"net/http"
)

func init() {
	// Load .env variables
	_ = gotenv.Load()
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}

