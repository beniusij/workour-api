package gql

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"net/http"
)

type reqBody struct {
	Query string `json:"query"`
}

// GraphQL returns an http.HandlerFunc to our /graphql endpoint
func GraphQL(sc graphql.Schema) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check to ensure query was provided in the request body
		if c.Request.Body == nil {
			http.Error(c.Writer, "Must provide graphql query in request body", http.StatusBadRequest)
			return
		}

		var rBody reqBody
		// Decode the request body into rBody
		err := json.NewDecoder(c.Request.Body).Decode(&rBody)
		if err != nil {
			http.Error(c.Writer, "Error parsing JSON request body", http.StatusBadRequest)
		}

		// Execute graphql query
		result, errs := ExecuteQuery(rBody.Query, sc, c)

		if len(errs) > 0 {
			c.JSON(http.StatusBadRequest, result)
			return
		}

		status, _ := c.Get("status")
		c.JSON(status.(int), result)
	}
}