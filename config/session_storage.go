package config

import (
	"fmt"
	redis "gopkg.in/boj/redistore.v1"
	"log"
)

var (
	store *redis.RediStore
	config *RedisConfig
)

func SetupSessionStorage() {
	config := NewRedisConfig()

	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	s, err := redis.NewRediStore(10, "tcp", address, config.Password, []byte(config.Secret))
	if err != nil {
		log.Panic(err)
	}

	store = s
}

func GetSessionStorage() *redis.RediStore {
	return store
}