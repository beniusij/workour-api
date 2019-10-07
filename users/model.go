package users

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"workour-api/common"
)

type User struct {
	ID				int		`gorm:"primary_key"`
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
	return bcrypt.CompareHashAndPassword(hashedPassword, bytePassword)
}

func (u User) SaveEntity(data interface{}) (int, error) {
	db := common.GetDB()
	user := data.(*User)
	err := db.Create(&user).Error

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (u User) GetEntityById(id int) (*User, error) {
	db := common.GetDB()
	user := User{}
	err := db.Where(&User{ID: id}).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}