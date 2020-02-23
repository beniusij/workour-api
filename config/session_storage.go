package config

import (
	"fmt"
	redis "gopkg.in/boj/redistore.v1"
	"os"
)

var store *redis.RediStore

// Initialise persistence cache with Redis
//func SetupSessionStorage(r *gin.Engine) {
//	address := fmt.Sprintf("%s:%s", host, port)
//	store, err := redis.NewStore(10, "tcp", address, psw, []byte(secret))
//	if err != nil {
//		panic(err)
//	}
//
//	r.Use(sessions.Sessions("auth_sessions", store))
//}

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
		host = "localhost"
	}

	// Port defaults to 6379
	if port = os.Getenv("REDIS_PORT"); port == "" {
		port = "6379"
	}

	// Default password can be empty for local dev env
	psw = os.Getenv("REDIS_PASSWORD")

	// Secret is whatever the fuck you want
	// Should not vary in environment
	if secret = os.Getenv("REDIS_SECRET"); secret == "" {
		secret = "VerySecureSecret"
	}
}