package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
)

const testDbPath = "./../gorm_test.db"
var DB *gorm.DB

func InitDb() *gorm.DB {
	var (
		driver = "postgres"
		dbc = NewDatabaseConfig()
	)

	creds := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		dbc.Host, dbc.Port, dbc.User, dbc.Name, dbc.Password, dbc.Secure)

	db, err := gorm.Open(driver, creds)

	if err != nil {
		panic(err)
	}

	db.DB().SetMaxIdleConns(10)

	DB = db
	return DB
}

func InitTestDb() * gorm.DB {
	var driver = "sqlite3"

	testDb, err := gorm.Open(driver, testDbPath)
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

func GetDB() *gorm.DB {
	return DB
}