package config

import (
	"fmt"
	redis "gopkg.in/boj/redistore.v1"
	"log"
	"os"
)

var store *redis.RediStore

var (
	host 	string
	port 	string
	psw 	string
	secret 	string
)

func SetupSessionStorage() {
	setEnvVars()

	address := fmt.Sprintf("%s:%s", host, port)
	s, err := redis.NewRediStore(10, "tcp", address, psw, []byte(secret))
	if err != nil {
		panic(err)
	}

	store = s
}

func GetSessionStorage() *redis.RediStore {
	return store
}

// Get env vars or set default values
func setEnvVars() {
	// Host defaults to localhost
	if host = os.Getenv("REDIS_HOST"); host == "" {
		log.Println("Setting default value for $REDIS_HOST")
		host = "localhost"
	}

	// Port defaults to 6379
	if port = os.Getenv("REDIS_PORT"); port == "" {
		log.Println("Setting default value for $REDIS_PORT")
		port = "6379"
	}

	// Default password can be empty for local dev env
	psw = os.Getenv("REDIS_PASSWORD")

	// Secret is whatever the fuck you want
	// Should not vary in environment
	if secret = os.Getenv("REDIS_SECRET"); secret == "" {
		log.Println("Setting default value for $REDIS_SECRET")
		secret = "VerySecureSecret"
	}
}