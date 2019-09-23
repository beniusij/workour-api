package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func initDb(driver, creds string) *gorm.DB {
	db, err := gorm.Open(driver, creds)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Migrate
	return db
}

func DBMiddleware(driver, creds string) gin.HandlerFunc {
	db := initDb(driver, creds)

	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}