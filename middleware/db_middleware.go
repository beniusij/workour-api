package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"workour-api/users"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

func initDb(driver, creds string) *gorm.DB {
	db, err := gorm.Open(driver, creds)

	if err != nil {
		panic(err)
	}

	db.DB().SetMaxIdleConns(10)

	// Migrate
	db.AutoMigrate(users.User{})

	DB = db
	return db
}

func DBMiddleware(driver, creds string) gin.HandlerFunc {
	db := initDb(driver, creds)

	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}

func GetDB() *gorm.DB {
	return DB
}