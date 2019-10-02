package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/subosito/gotenv"
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

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&u.User{})
}

func main() {
	r, db := initAPI()
	Migrate(db)
	defer db.Close()

	r.Run(":8080")
}

func initAPI() (*gin.Engine, *gorm.DB) {
	db := comm.InitDb()
	router := gin.Default()

	rootQuery := g.NewRoot()
	// Create a new graphql schema, passing in the root query
	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query: rootQuery.Query,
		},
	)
	if err != nil {
		fmt.Println("error creating schema: ", err)
	}

	router.POST("/graphql", g.GraphQL(schema))

	return router, db
}