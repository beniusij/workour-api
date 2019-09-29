package gql

import (
	"encoding/json"
	"github.com/gin-gonic/gin/render"
	"github.com/graphql-go/graphql"
	"net/http"
)

type reqBody struct {
	Query string `json:"query"`
}

// GraphQL returns an http.HandlerFunc to our /graphql endpoint
func GraphQL(sc graphql.Schema) http.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check to ensure query was provided in the request body
		if r.Body == nil {
			http.Error(w, "must provide graphql query in request body", http.StatusUnprocessableEntity)
			return
		}
		var rBody reqBody
		// Decode the request body into rBody
		err := json.NewDecoder(r.Body).Decode(&rBody)
		if err != nil {
			http.Error(w, "error parsing JSON request body", http.StatusUnprocessableEntity)
		}

		// Execute graphql query
		result := gql.ExecuteQuery(rBody.Query, sc)

		//render.JSON(w, r, result)
	}
}