package tests

import (
	"database/sql"
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
	router := gin.New()

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
	roleMocker()

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

func DeleteCreatedEntities(db *gorm.DB) func() {
	type entity struct {
		table   string
		keyname string
		key     interface{}
	}
	var entries []entity
	hookName := "cleanupHook"

	// Setup the onCreate Hook
	db.Callback().Create().After("gorm:create").Register(hookName, func(scope *gorm.Scope) {
		fmt.Printf("Inserted entities of %s with %s=%v\n", scope.TableName(), scope.PrimaryKey(), scope.PrimaryKeyValue())
		entries = append(entries, entity{table: scope.TableName(), keyname: scope.PrimaryKey(), key: scope.PrimaryKeyValue()})
	})
	return func() {
		// Remove the hook once we're done
		defer db.Callback().Create().Remove(hookName)
		// Find out if the current db object is already a transaction
		_, inTransaction := db.CommonDB().(*sql.Tx)
		tx := db
		if !inTransaction {
			tx = db.Begin()
		}
		// Loop from the end. It is important that we delete the entries in the
		// reverse order of their insertion
		for i := len(entries) - 1; i >= 0; i-- {
			entry := entries[i]
			fmt.Printf("Deleting entities from '%s' table with key %v\n", entry.table, entry.key)
			tx.Table(entry.table).Where(entry.keyname+" = ?", entry.key).Delete("")
		}

		if !inTransaction {
			tx.Commit()
		}
	}
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

func roleMocker() {
	role := roles.Role{
		Name:      "Regular User",
		Authority: 1,
		Policies:  nil,
	}
	db.Create(&role)
}