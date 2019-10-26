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

	// Verify that generated string matches JWT token regex
	token := comm.GenerateToken(creds["email"].(string))
	asserts.NotNil(token, "token is not nill")
	asserts.Regexp(tokenRegex, token, "JWT token matches token regex")

	// Verify that the JWT contains three segments, separated by two period ('.') characters
	tokenSlice := strings.Split(token, ".")
	asserts.EqualValues(3, len(tokenSlice), "JWT token contains 3 parts separated by '.'")
}