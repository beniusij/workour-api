package common

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
	"time"
)

type CommonError struct {
	Errors map[string]interface{} `json:"errors"`
}

func Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(obj, b)
}

func NewError(key string, err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	res.Errors[key] = err.Error()
	return res
}

func NewValidationError(err error) CommonError {
	res := CommonError{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)

	for _, v := range errs {
		if v.Param != "" {
			res.Errors[v.Field] = fmt.Sprintf("{%v: %v}", v.Tag, v.Param)
		} else {
			res.Errors[v.Field] = fmt.Sprintf("{key: %v}", v.Tag)
		}
	}

	return res
}

// This one should be private
const NBSecretPassword = "This is 4 bl00dy s3cur3 p455w0rd!#"

func GetToken(id int) string {
	jwtToken := jwt.New(jwt.GetSigningMethod("HS256"))
	// Set some claims
	jwtToken.Claims = jwt.MapClaims{
		"id":	id,
		"exp":	time.Now().Add(time.Hour * 24).Unix(),
	}
	// Sign and get the complete encoded token as a string
	token, _ := jwtToken.SignedString([]byte(NBSecretPassword))
	return token
}