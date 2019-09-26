package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"workour-api/helpers"
)

func CreateUser(c *gin.Context) {
	modelValidator := UserModelValidator{}
	err := modelValidator.Bind(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, helpers.NewValidationError(err))
	}

	err = SaveUser(&modelValidator.user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, helpers.NewError("database", err))
		return
	}


}