package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/subosito/gotenv"
	"os"
	comm "workour-api/common"
	"workour-api/config"
	u "workour-api/users"
)

var (
	appPort = os.Getenv("GIN_PORT")

)

func init() {
	// Load .env variables
	err := gotenv.Load()

	if err != nil {
		panic(err)
	}
}

func main() {
	r, db := initAPI()
	Migrate(db)
	defer db.Close()

	_ = r.Run(":8080")
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		u.User{},
	)
}

func initAPI() (*gin.Engine, *gorm.DB) {
	db := comm.InitDb()
	router := gin.Default()

	// Set up Redis client for sessions
	comm.InitSessionStore(router)

	config.SetupRoutes(router)

	return router, db
}

