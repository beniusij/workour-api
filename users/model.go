package users

import (
	//"errors"
	"errors"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"workour-api/middleware"
)

type User struct {
	gorm.Model
	Email			string	`gorm:"column:username;type:varchar(100);unique_index"`
	FirstName		string	`gorm:"column:first_name"`
	LastName		string	`gorm:"column:last_name"`
	PasswordHash	string	`gorm:"column:password;not null"`
}

func SaveUser(data interface{}) error {
	db := middleware.GetDB()
	err := db.Save(data).Error
	return err
}

func (u *User) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password cannot be empty")
	}

	bytePassword := []byte(password)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.PasswordHash = string(passwordHash)

	return nil
}