package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"workour-api/common"
)

func CreateUser(c *gin.Context) {
	modelValidator := UserModelValidator{}
	err := modelValidator.Bind(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidationError(err))
	}

	err = SaveUser(&modelValidator.user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}


}