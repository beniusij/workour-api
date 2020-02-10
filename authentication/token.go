package common

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
	"workour-api/users"
)

type JWTToken interface {
	GenerateToken(email string) (string, error)
	ValidateToken() bool
	DecodeToken() (map[string]interface{}, error)
	RefreshToken() (string, error)
}

type AuthToken struct {
	Token string
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func (t AuthToken)GenerateToken(u users.User) (string, error) {
	claims := jwt.MapClaims{
		"id": u.ID,
		"email": u.Email,
		"fname": u.FirstName,
		"lname": u.LastName,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var err error
	t.Token, err = token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return t.Token, nil
}

func (t *AuthToken)ValidateToken() bool {
	// Verify that signed Token is not empty
	if t.Token == "" {
		return false
	}

	// Get unsigned Token and verify no errors occurred in the process
	token, err := jwt.Parse(t.Token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecret, nil
	})
	if err != nil {
		return false
	}

	// verify Token as not expired
	claims := token.Claims.(jwt.MapClaims)
	return claims.VerifyExpiresAt(jwt.TimeFunc().Unix(), false)
}

func (t *AuthToken)DecodeToken() (map[string]interface{}, error) {
	// Get unsigned Token
	token, err := jwt.Parse(t.Token, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	// Map claims attributes from object to map structure
	claims := token.Claims.(jwt.MapClaims)
	claimsMap := make(map[string]interface{})
	for key, val := range claims {
		claimsMap[key] = val
	}

	return claimsMap, nil
}

func (t *AuthToken)RefreshToken() (string, error) {
	// extract claims from signed Token
	token, err := jwt.Parse(t.Token, func(token *jwt.Token) (i interface{}, e error) {
		return jwtSecret, nil
	})
	if err != nil {
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)

	// set new expiry date
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// generate Token with updated claims
	refreshedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t.Token, err = refreshedToken.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return t.Token, nil
}