package main

import (
	"github.com/jinzhu/gorm"
	"github.com/subosito/gotenv"
	m "workour-api/middleware"
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
	db := m.InitDb()
	Migrate(db)
	defer db.Close()

	r := SetupRouter()
	r.Run(":8080")
}

