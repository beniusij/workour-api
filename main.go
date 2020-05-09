package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"workour-api/config"
	r "workour-api/roles"
	u "workour-api/users"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found, will load runtime env variables")
	}
}

func main() {
	config.New()
	r, db := initAPI()

	Migrate(db)
	defer db.Close()

	_ = r.Run(fmt.Sprintf(":%s", config.Configurations.Port))
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

