package models

import (
	//"errors"
	"errors"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Email			string	`gorm:"column:username;type:varchar(100);unique_index"`
	FirstName		string	`gorm:"column:first_name"`
	LastName		string	`gorm:"column:last_name"`
	PasswordHash	string	`gorm:"column:password;not null"`
}

//func (u *User) New(data interface{}) (bool) {
//	//db, _ := c.Context{}.Get("db")
//	//user := User{}
//}

func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password cannot be empty")
	}

	bytePassword := []byte(password)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.PasswordHash = string(passwordHash)

	return nil
}