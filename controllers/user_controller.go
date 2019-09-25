package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"workour-api/helpers"
	"workour-api/validators"
)

func CreateUser(c *gin.Context) {
	modelValidator := validators.UserModelValidator{}
	err := modelValidator.Bind(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, helpers.NewValidationError(err))
	}

	
}