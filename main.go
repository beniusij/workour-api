package main

import (
	"fmt"
	"github.com/subosito/gotenv"
	"workour-api/config"
	"os"
)

func init() {
	// Load .env variables
	err := gotenv.Load()

	if err != nil {
		panic(err)
	}
}

func main() {
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

	r := config.SetupRouter(driver, creds)
	r.Run(":8080")
}

