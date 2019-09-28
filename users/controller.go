package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"workour-api/common"
)

func CreateUser(c *gin.Context) {
	modelValidator := NewUserModelValidator()
	err := modelValidator.Bind(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidationError(err))
		return
	}

	err = SaveUser(&modelValidator.userModel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	c.Set("my_user_model", modelValidator.userModel)
	serializer := UserSerializer{c}
	c.JSON(http.StatusCreated, gin.H{"user": serializer.Response()})
}