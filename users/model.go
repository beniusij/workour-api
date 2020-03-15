package users

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"workour-api/config"
	"workour-api/roles"
)

type User struct {
	ID				uint		`gorm:"primary_key"`
	Role			roles.Role 	`gorm:"foreignkey:RoleId"`
	RoleId			uint		`gorm:"not null"`
	Email			string		`gorm:"unique"`
	FirstName		string
	LastName		string
	PasswordHash	string		`gorm:"column:password;not null"`
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

func (u User) Save(data interface{}) (User, error) {
	db := config.GetDB()
	user := data.(User)

	// Check and get default role
	if user.RoleId == uint(0) {
		user.RoleId	= roles.GetDefaultRoleId()
	}

	err := db.Create(&user).Error
	if err != nil {
		return User{}, err
	}

	return user, nil
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