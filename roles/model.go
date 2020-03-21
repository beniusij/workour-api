package roles

import (
	"github.com/jinzhu/gorm"
	"log"
	"reflect"
	"workour-api/config"
)

const defaultRoleId = "Regular User"

type Role struct {
	gorm.Model
 	Name 		string `gorm:"unique;not null"`
	Authority 	int
	Policies 	[]Policy `gorm:"foreignkey:RoleId"`
}

func GetDefaultRoleId() uint {
	db := config.GetDB()
	role := Role{Name: defaultRoleId}

	err := db.First(&role).Error
	if err != nil {
		log.Fatalln(err)
	}

	return role.ID
}

// Get role by its ID
func (r* Role) GetById() error {
	db := config.GetDB()
	err := db.Where(&r).First(&r).Error

	return err
}

// Get Role and its Policies by role id
func GetRoleById(id uint) Role {
	db := config.GetDB()
	role := Role{}

	role.ID = id
	db.First(&role)

	return role
}

type Policy struct {
	gorm.Model
	RoleId		uint
	Resource 	string `gorm:"not null"`
	Index		bool
	Create 		bool
	Read 		bool
	Update 		bool
	Delete 		bool
}

// Gets field value by field name
func (p* Policy) GetFieldValueByName(fieldName string) bool {
	r := reflect.ValueOf(p)
	a := reflect.Indirect(r).FieldByName(fieldName)
	return a.Bool()
}