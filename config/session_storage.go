package config

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"os"
)

var (
	host = os.Getenv("REDIS_HOST")
	port = os.Getenv("REDIS_PORT")
	psw = os.Getenv("REDIS_PASSWORD")
	secret = os.Getenv("REDIS_SECRET")
)

// Initialise persistence cache with Redis
func SetupSessionStorage(r *gin.Engine) {
	address := fmt.Sprintf("%s:%s", host, port)
	store, err := redis.NewStore(10, "tcp", address, psw, []byte(secret))
	if err != nil {
		panic(err)
	}

	r.Use(sessions.Sessions("auth_sessions", store))
}