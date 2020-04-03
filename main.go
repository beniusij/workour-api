package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"workour-api/config"
	r "workour-api/roles"
	u "workour-api/users"
)

var PORT = os.Getenv("PORT")

func main() {
	if PORT != "" {
		log.Fatal("$PORT must be set")
	}

	r, db := initAPI()
	Migrate(db)
	defer db.Close()

	_ = r.Run(fmt.Sprintf(":%s", PORT))
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		u.User{},
		r.Policy{},
		r.Role{},
	)

	r.CreateDefaultRoles()
}

func initAPI() (*gin.Engine, *gorm.DB) {
	db := config.InitDb()
	router := gin.Default()

	// Set up Redis client for sessions
	config.SetupSessionStorage()

	SetupRoutes(router)

	return router, db
}

