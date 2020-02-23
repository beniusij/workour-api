package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	auth "workour-api/authentication"
	g "workour-api/gql"
)

func SetupRoutes(router *gin.Engine) {
	publicRoutes(router)
	adminRoutes(router)
}

// Set up public level routes
func publicRoutes(r *gin.Engine) {
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

	r.POST("/login", authController.AuthenticateUser)
	r.POST("/logout", authController.LogoutUser)
	r.POST("/public", g.GraphQL(schema))
	r.OPTIONS("/public", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Request-Method","POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept-Encoding")
	})
}

// Set up role-protected routes
func adminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	admin.Use(auth.VerifyAuthentication())
}