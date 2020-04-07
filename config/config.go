package config

import (
	"os"
	"strings"
)

// Database configuration struct
type DatabaseConfig struct {
	Host 		string
	Port		string
	User 		string
	Password 	string
	Name		string
	Secure		string
}

// Redis configuration struct
type RedisConfig struct {
	Host		string
	Port		string
	Password	string
	Secret		string
}

type DatabasesConfig struct {
	db 		DatabaseConfig
	tdb 	DatabaseConfig
}

// Configs loaded from environment
type Config struct {
	TestDatabase 	DatabaseConfig
	Port			string
	JWTSecret		string
	CorsOrigins		[]string
	Environment		string
	Domain			string
}

// Returns a set of configs
func New() *Config {
	return &Config{
		Port: 			getEnv("PORT", "8080"),
		JWTSecret:		getEnv("JWT_SECRET", "ThisIsTokenSecret"),
		CorsOrigins:	getEnvAsSlice("CORS_ORIGIN", []string{"http://localhost:3000"}, ","),
		Environment:	getEnv("APP_ENV", "development"),
		Domain:			getEnv("DOMAIN", "localhost"),
	}
}

func NewRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		Secret:   getEnv("REDIS_SECRET", "NothingSoSecretHere"),
	}
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DATABASE_HOST", "workour_db"),
		Port:     getEnv("DATABASE_PORT", "5432"),
		User:     getEnv("DATABASE_USER", "root"),
		Password: getEnv("DATABASE_PSW", "admin"),
		Name:     getEnv("DATABASE_NAME", "workour"),
		Secure:   getEnv("DATABASE_SSL", "disable"),
	}
}

// Simple helper function to read an environment or return a default value
// Reference: https://dev.to/craicoverflow/a-no-nonsense-guide-to-environment-variables-in-go-a2f
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a string slice or return default value
// Reference: https://dev.to/craicoverflow/a-no-nonsense-guide-to-environment-variables-in-go-a2f
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}