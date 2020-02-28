package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"os"
	"time"
	auth "workour-api/authentication"
	g "workour-api/gql"
)

var allowHeaders = []string{
	"Accept",
	"Accept-Encoding",
	"Authorization",
	"Content-Length",
	"Content-Type",
	"X-CSRF-Token",
}

func SetupRoutes(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:           []string{os.Getenv("CORS_ORIGIN")},
		AllowMethods:           []string{"POST", "PUT", "GET", "OPTIONS"},
		AllowHeaders:           allowHeaders,
		AllowCredentials:       true,
		ExposeHeaders:          []string{"Content-Length"},
		MaxAge:                 24 * time.Hour,
		AllowFiles:             true,
	}))

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
	r.GET("/getCurrentUser", authController.GetCurrentUser)
	r.POST("/public", g.GraphQL(schema))
}

// Set up role-protected routes
func adminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	admin.Use(auth.VerifyAuthentication())
}