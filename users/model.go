package users

import (
	//"errors"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"workour-api/common"
)

type User struct {
	ID				uint	`gorm:"primary_key"`
	Email			string	`gorm:"column:email;type:varchar(100);unique_index"`
	FirstName		string	`gorm:"column:first_name"`
	LastName		string	`gorm:"column:last_name"`
	PasswordHash	string	`gorm:"column:password;not null"`
}

func AutoMigrate() {
	db := common.GetDB()

	db.AutoMigrate(&User{})
}

func SaveUser(data interface{}) error {
	db := common.GetDB()
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

/**
 * Checks passed password against the hashed password
 * using bcrypt's function
 */
func (u *User) CheckPassword(password string) error {
	bytePassword := []byte(password)
	hashedPassword := []byte(u.PasswordHash)
	return bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
}

func GetUserById(id uint) (User, error) {
	db := common.GetDB()
	var model User
	err := db.Where(&User{ID: id}).First(&model).Error
	return model, err
}