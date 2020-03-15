package roles

import (
	"github.com/jinzhu/gorm"
)

type Policy struct {
	gorm.Model
	Resource 	string `gorm:"unique;not null"`
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