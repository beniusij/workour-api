package common

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type JWTToken interface {
	GenerateToken(email string) (string, error)
	ValidateToken() bool
	DecodeToken() (map[string]interface{}, error)
	RefreshToken() (string, error)
}

type AuthToken struct {
	token string
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func (t AuthToken)GenerateToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var err error
	t.token, err = token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return t.token, nil
}

func (t *AuthToken)ValidateToken() bool {
	// Verify that signed token is not empty
	if t.token == "" {
		return false
	}

	// Get unsigned token and verify no errors occured in the process
	token, err := jwt.Parse(t.token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecret, nil
	})
	if err != nil {
		return false
	}

	// verify token as not expired
	claims := token.Claims.(jwt.MapClaims)
	return claims.VerifyExpiresAt(jwt.TimeFunc().Unix(), false)
}

func (t *AuthToken)DecodeToken(string) (map[string]interface{}, error) {
	// Get unsigned token
	token, err := jwt.Parse(t.token, func(token *jwt.Token) (interface{}, error) {
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
	// extract claims from signed token
	token, err := jwt.Parse(t.token, func(token *jwt.Token) (i interface{}, e error) {
		return jwtSecret, nil
	})
	if err != nil {
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)

	// set new expiry date
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// generate token with updated claims
	refreshedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t.token, err = refreshedToken.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return t.token, nil
}