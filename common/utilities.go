package common

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

// This one should be private
const NBSecretPassword = "This is 4 bl00dy s3cur3 p455w0rd!#"

func GetToken(id uint) string {
	jwtToken := jwt.New(jwt.GetSigningMethod("HS256"))
	// Set some claims
	jwtToken.Claims = jwt.MapClaims{
		"id":	id,
		"exp":	time.Now().Add(time.Hour * 24).Unix(),
	}
	// Sign and get the complete encoded Token as a string
	token, _ := jwtToken.SignedString([]byte(NBSecretPassword))
	return token
}