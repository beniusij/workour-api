package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/subosito/gotenv"
	"os"
	auth "workour-api/authentication"
	comm "workour-api/common"
	g "workour-api/gql"
	u "workour-api/users"
)

func init() {
	// Load .env variables
	err := gotenv.Load()

	if err != nil {
		panic(err)
	}
}

func main() {
	r, db := initAPI()
	Migrate(db)
	defer db.Close()

	_ = r.Run(":8080")
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		u.User{},
	)
}

func initAPI() (*gin.Engine, *gorm.DB) {
	db := comm.InitDb()
	router := gin.Default()

	// Set up Redis client for sessions
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	address := fmt.Sprintf("%s:%s", redisHost, redisPort)
	comm.InitSessionStore(router, address)

	rootQuery := g.NewRoot()
	// Create a new graphql schema, passing in the root query
	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: rootQuery.Query,
			Mutation: rootQuery.Mutation,
		},
	)

	if err != nil {
		fmt.Println("error creating schema: ", err)
	}

	authController := new(auth.Controller)

	// Set up public level routes
	router.POST("/login", authController.AuthenticateUser)
	router.POST("/logout", authController.LogoutUser)
	router.POST("/public", g.GraphQL(schema))
	router.OPTIONS("/public", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Request-Method","POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept-Encoding")
	})

	// Set up role-protected routes
	admin := router.Group("/admin")
	admin.Use(auth.VerifyAuthentication())

	return router, db
}

