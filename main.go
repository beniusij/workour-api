package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/subosito/gotenv"
	"workour-api/common"
	c "workour-api/config"
	"workour-api/gql"
	u "workour-api/users"
)

func init() {
	// Load .env variables
	err := gotenv.Load()

	if err != nil {
		panic(err)
	}
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&u.User{})
}

func main() {
	r, db := initAPI()
	defer db.Close()

	r.Run(":8080")
}

func initAPI() (*gin.Engine, *gorm.DB) {
	db := common.InitDb()
	Migrate(db)

	router := c.SetupRouter()

	rootQuery := gql.NewRoot(db)
	

	return router, db
}