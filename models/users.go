package models

import (
	//"errors"
	"github.com/jinzhu/gorm"
	//"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Email			string	`gorm:"column:username;type:varchar(100);unique_index"`
	FirstName		string	`gorm:"column:first_name"`
	LastName		string	`gorm:"column:last_name"`
	PasswordHash	string	`gorm:"column:password;not null"`
}