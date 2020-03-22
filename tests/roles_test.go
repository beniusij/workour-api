package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"workour-api/roles"
)

func TestGetRoleAndRelatedPolicies(t *testing.T) {
	// Set up cleaner hook
	cleaner := deleteCreatedEntities(db)
	defer cleaner()

	// Set up test fixtures
	roleMocker()

	asserts := assert.New(t)

	t.Run("Should get role by ID and associated policies", func(t *testing.T) {
		role := roles.Role{}
		role.ID = roles.GetDefaultRoleId()

		err := role.GetById()

		asserts.NoError(err, "Should not return an error")
		asserts.Equal(defaultRole, role.Name, fmt.Sprintf("Role name should be %s", defaultRole))
		asserts.NotNil(role.Policies, "Should have policies")
	})

	t.Run("Should return error and blank role struct", func(t *testing.T) {
		role := roles.Role{}
		role.ID = 100

		err := role.GetById()
		t.Log(err)

		asserts.Error(err, "Should return record not found error")
		asserts.Equal("", role.Name, "Should return blank role struct")
		asserts.Equal(0, len(role.Policies), "Should not have policies")
	})
}
