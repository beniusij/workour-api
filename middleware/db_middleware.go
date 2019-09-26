package middleware

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

func InitDb() *gorm.DB {
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

	db, err := gorm.Open(driver, creds)

	if err != nil {
		panic(err)
	}

	db.DB().SetMaxIdleConns(10)

	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}