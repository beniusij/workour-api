package tests

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"testing"
	"time"
)

type StubAuthToken struct {
	token string
}

var creds = map[string]interface{}{
	"email":	"userModel1@yahoo.com",
	"password":	"Password123",
}

var jwtSecret = []byte("SecretStringForSigningTokens")
var TokenRegex = `^[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*$`

func (t *StubAuthToken)GenerateToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var err error
	t.token, err = token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return t.token, nil
}

func (t *StubAuthToken)ValidateToken() bool {
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

func (t *StubAuthToken)DecodeToken() (map[string]interface{}, error) {
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

func (t *StubAuthToken)RefreshToken() (string, error) {
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

func initTestToken() (StubAuthToken, error) {
	var token StubAuthToken
	tokenString, err := token.GenerateToken(creds["email"].(string))
	if err != nil {
		return StubAuthToken{token:""}, err
	}
	token.token = tokenString
	return token, nil
}

// ---------------------------------------------------------------------------
// ---------------------------------- Tests ----------------------------------
// ---------------------------------------------------------------------------

func TestGenerateToken(t *testing.T) {
	asserts := getAsserts(t)
	token, err := initTestToken()

	// Verify that generated string matches JWT token regex
	asserts.NoError(err, "no error is returned")
	asserts.NotNil(token, "token is not nil")
	asserts.Regexp(TokenRegex, token.token, "JWT token matches token regex")

	// Verify that the JWT contains three segments, separated by two period ('.') characters
	tokenSlice := strings.Split(token.token, ".")
	asserts.EqualValues(3, len(tokenSlice), "JWT token contains 3 parts separated by '.'")
}

func TestValidateToken(t *testing.T) {
	asserts := getAsserts(t)
	token, err := initTestToken()

	asserts.NoError(err, "should not return error")

	validToken := token.ValidateToken()
	token.token = "Potato123"
	invalidToken := token.ValidateToken()

	asserts.True(validToken, "should be valid")
	asserts.False(invalidToken, "should be invalid")
}

func TestDecodeToken(t *testing.T) {
	asserts := getAsserts(t)
	token, _ := initTestToken()

	t.Run("valid token should be decoable and return map with email and exp", func(t *testing.T) {
		claimsMap, err := token.DecodeToken()

		asserts.NoError(err, "no error returned while decoding token")
		asserts.IsType(make(map[string]interface{}), claimsMap, "claims map is of type map[string]interface{}")
		asserts.EqualValues(creds["email"].(string), claimsMap["email"].(string), "should return right email")
	})

	t.Run("invalid token should only return error and no claims map", func(t *testing.T) {
		tokens := []string{
			"",
			"totally.invalid.string",
			"123456!+=",
		}

		for _, t := range tokens {
			token.token = t
			claimsMap, err := token.DecodeToken()

			asserts.Nil(claimsMap, "claims map should be nil")
			asserts.Error(err, "should return error")
		}
	})
}

func TestRefreshToken(t *testing.T) {
	asserts := getAsserts(t)
	token, _ := initTestToken()

	oldToken := token.token
	claims, _ := token.DecodeToken()
	refreshedToken, err := token.RefreshToken()
	refreshedClaims, _ := token.DecodeToken()
	oldExp := time.Unix(int64(claims["exp"].(float64)), 0)
	newExp := time.Unix(int64(refreshedClaims["exp"].(float64)), 0)

	// Refreshed token is valid and is different than the original
	asserts.NoError(err, "should not return an error")
	asserts.NotNil(refreshedToken, "should not be nil")
	asserts.True(token.ValidateToken(), "should be valid")

	// Refresh token claims have the same email, expires later than the original token
	asserts.True(oldToken != refreshedToken, "should not be the same token")
	asserts.EqualValues(creds["email"].(string), refreshedClaims["email"].(string), "should have the same email")
	asserts.True(oldExp.Before(newExp), "should expire later")
}