package tests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
	"workour-api/config"
	"workour-api/gql"
	"workour-api/roles"
	u "workour-api/users"
)

const publicEndpoint = "/public"
var db *gorm.DB
var asserts *assert.Assertions

// ------------------------------------------------------------------------------------
// ------------------------------- Tools initialisation -------------------------------
// ------------------------------------------------------------------------------------

func initTestAPI() *gin.Engine {
	db = config.InitTestDb()
	migrate()
	router := gin.Default()

	rootQuery := gql.NewRoot()
	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:		rootQuery.Query,
			Mutation:   rootQuery.Mutation,
		},
	)
	if err != nil {
		panic(err)
	}

	router.POST(publicEndpoint, gql.GraphQL(schema))

	return router
}

func getAsserts(t *testing.T) *assert.Assertions {
	if asserts == nil {
		asserts = assert.New(t)
	}

	return asserts
}

// ------------------------------------------------------------------------------
// ------------------------------- Database utils -------------------------------
// ------------------------------------------------------------------------------

func resetDb(addMock bool) {
	_ = config.ResetTestDb(db)
	db = config.InitTestDb()
	migrate()

	addMockRoles()

	if addMock {
		userMocker(10)
	}
}

func migrate() {
	db.AutoMigrate(
		u.User{},
		roles.Policy{},
		roles.Role{},
	)
}

// -------------------------------------------------------------------------
// ------------------------------- Mock ------------------------------------
// -------------------------------------------------------------------------

func userMocker(n int) []u.User {
	var offset int
	var ret []u.User
	db.Model(&u.User{}).Count(&offset)

	for i := offset + 1; i <= offset+n; i++ {
		user := u.User{
			Email: fmt.Sprintf("userModel%v@yahoo.com", i),
			FirstName: fmt.Sprintf("User%v", i),
			LastName: fmt.Sprintf("User%v", i),
			RoleId: roles.GetDefaultRoleId(),
		}
		_ = user.SetPassword("Password123")
		db.Create(&user)
		ret = append(ret, user)
	}

	return ret
}

func addMockRoles() {
	role := roles.Role{
		Name:      "Regular User",
		Authority: 1,
		Policies:  nil,
	}
	db.Create(&role)
}

func newUserModel() u.User {
	return u.User{
		Email: "t3st@gmail.com",
		FirstName: "Testas",
		LastName: "Testavicius",
		PasswordHash: "",
	}
}