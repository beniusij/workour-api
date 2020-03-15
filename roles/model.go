package roles

import (
	"github.com/jinzhu/gorm"
	"log"
	"workour-api/config"
)

const defaultRoleId = "Regular User"

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