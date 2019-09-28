package main

import (
	"github.com/jinzhu/gorm"
	"github.com/subosito/gotenv"
	"workour-api/common"
	c "workour-api/config"
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
	db := common.InitDb()
	Migrate(db)
	defer db.Close()

	r := c.SetupRouter()
	r.Run(":8080")
}

