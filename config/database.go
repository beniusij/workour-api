package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
)

var DB *gorm.DB
var testDbPath = "./../gorm_test.db"

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
		fmt.Println(err)
	}

	db.DB().SetMaxIdleConns(10)

	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}

func InitTestDb() * gorm.DB {
	testDb, err := gorm.Open("sqlite3", testDbPath)
	if err != nil {
		panic(err)
	}

	testDb.DB().SetMaxIdleConns(3)
	testDb.LogMode(false)
	DB = testDb
	return testDb
}

func ResetTestDb(db *gorm.DB) error {
	db.Close()
	err := os.Remove(testDbPath)
	return err
}