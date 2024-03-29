package gql

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"net/http"
)

type reqBody struct {
	Query string `json:"query"`
	Variables map[string]interface{} `json:"variables"`
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
		result, errs := ExecuteQuery(rBody.Query, rBody.Variables, sc, c)

		if len(errs) > 0 {
			c.JSON(http.StatusOK, result)
			return
		}

		status, _ := c.Get("status")
		if status == nil {
			fmt.Println(fmt.Sprint("no status is returned"))
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		c.JSON(status.(int), result)
	}
}