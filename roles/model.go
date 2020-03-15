package roles

import (
	"github.com/jinzhu/gorm"
)

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