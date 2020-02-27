package users

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"workour-api/config"
)

type User struct {
	ID				uint	`gorm:"primary_key"`
	Email			string	`gorm:"column:email;type:varchar(100);unique_index"`
	FirstName		string	`gorm:"column:first_name"`
	LastName		string	`gorm:"column:last_name"`
	PasswordHash	string	`gorm:"column:password;not null"`
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


// Checks passed password against the hashed password
// using bcrypt's function
func (u *User) CheckPassword(password string) error {
	bytePassword := []byte(password)
	hashedPassword := []byte(u.PasswordHash)

	if bcrypt.CompareHashAndPassword(hashedPassword, bytePassword) != nil {
		return errors.New("incorrect email and/or password")
	}

	return nil
}

func (u User) Save(data interface{}) (uint, error) {
	db := config.GetDB()
	user := data.(User)
	err := db.Create(&user).Error

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (u User) GetById(id uint) (User, error) {
	db := config.GetDB()
	user := User{}
	err := db.Where(&User{ID: id}).First(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func GetByEmail(email string) (User, error) {
	db := config.GetDB()
	user := User{}

	err := db.Where("email = ?", email).Find(&user).Error

	return user, err
}