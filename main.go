package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
	"net/http"
	"os"
	"workour-api/middleware"
)

func init() {
	// Load .env variables
	err := gotenv.Load()

	if err != nil {
		panic(err)
	}
}

func setupRouter(driver, creds string) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.DBMiddleware(driver, creds))

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}

func main() {
	var (
		driver = "postgres"
		dbUser = os.Getenv("DATABASE_USER")
		dbPsw = os.Getenv("DATABASE_PSW")
		dbHost = os.Getenv("DATABASE_HOST")
		dbPort = os.Getenv("DATABASE_PORT")
		dbName = os.Getenv("DATABASE_NAME")
		dbSSL = os.Getenv("DATABASE_SSL")
	)

	creds := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbName, dbPsw, dbSSL)

	r := setupRouter(driver, creds)
	r.Run(":8080")
}

