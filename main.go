package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"os"
	"workour-api/config"
	r "workour-api/roles"
	u "workour-api/users"
)

func main() {
	r, db := initAPI()
	Migrate(db)
	defer db.Close()

	_ = r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
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

