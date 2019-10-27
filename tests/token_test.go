package tests

import (
	"strings"
	"testing"
	comm "workour-api/common"
)

var creds = map[string]interface{}{
	"email":	"userModel1@yahoo.com",
	"password":	"Password123",
}

var tokenRegex = `^[A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*$`

func TestGenerateToken(t *testing.T) {
	asserts := getAsserts(t)
	token, err := comm.GenerateToken(creds["email"].(string))

	// Verify that generated string matches JWT token regex
	asserts.NoError(err, "no error is returned")
	asserts.NotNil(token, "token is not nil")
	asserts.Regexp(tokenRegex, token, "JWT token matches token regex")

	// Verify that the JWT contains three segments, separated by two period ('.') characters
	tokenSlice := strings.Split(token, ".")
	asserts.EqualValues(3, len(tokenSlice), "JWT token contains 3 parts separated by '.'")
}

func TestValidateToken(t *testing.T) {
	asserts := getAsserts(t)
	token, _ := comm.GenerateToken(creds["email"].(string))

	validToken := comm.ValidateToken(token)
	invalidToken := comm.ValidateToken("Potato123")

	asserts.True(validToken, "token is valid")
	asserts.False(invalidToken, "token is invalid")
}

func TestDecodeToken(t *testing.T) {
	asserts := getAsserts(t)
	tokenString, _ := comm.GenerateToken(creds["email"].(string))

	t.Run("valid token should be decoable and return map with email and exp", func(t *testing.T) {
		claimsMap, err := comm.DecodeToken(tokenString)

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

		for _, token := range tokens {
			claimsMap, err := comm.DecodeToken(token)

			asserts.Nil(claimsMap, "claims map should be nil")
			asserts.Error(err, "should return error")
		}
	})
}

func TestRefreshToken(t *testing.T) {
	asserts := getAsserts(t)
	tokenString, _ := comm.GenerateToken(creds["email"].(string))
	refreshedTokenString, err := comm.RefreshToken(tokenString)

	// Refreshed token is valid and is different than the original
	asserts.NoError(err, "should not return an error")
	asserts.NotNil(refreshedTokenString, "should not be nil")
	asserts.True(comm.ValidateToken(refreshedTokenString), "should be valid")
	asserts.True(tokenString != refreshedTokenString, "should not be the same token")

	claims, _ := comm.DecodeToken(tokenString)
	refreshedClaims, _ := comm.DecodeToken(refreshedTokenString)

	// Refresh token claims have the same email, expires later than the original token
	asserts.EqualValues(creds["email"].(string), refreshedClaims["email"].(string), "should have the same email")
	asserts.True(refreshedClaims["exp"].(float64) > claims["exp"].(float64), "should expire later")
}