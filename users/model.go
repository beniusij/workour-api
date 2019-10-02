package users

import (
	"errors"
	"fmt"
	g "github.com/graphql-go/graphql"
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

var UserType = g.NewObject(
	g.ObjectConfig{
		Name:        "User",
		Fields:      g.Fields{
			"ID": &g.Field{
				Type: g.Int,
			},
			"Email": &g.Field{
				Type: g.String,
			},
			"FirstName": &g.Field{
				Type: g.String,
			},
			"LastName": &g.Field{
				Type: g.String,
			},
		},
	},
)

var UserRegisterType = g.NewObject(
	g.ObjectConfig{
		Name:        "UserRegister",
		Fields:      g.Fields{
			"ID": &g.Field{
				Type: g.Int,
			},
			"Email": &g.Field{
				Type: g.String,
			},
			"FirstName": &g.Field{
				Type: g.String,
			},
			"LastName": &g.Field{
				Type: g.String,
			},
			"Password": &g.Field{
				Type: g.String,
			},
			"PasswordConfirm": &g.Field{
				Type: g.String,
			},
		},
	},
)

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

func GetUserById(id int) (User, error) {
	db := common.GetDB()
	var model User
	err := db.Where(&User{ID: id}).First(&model).Error

	if err != nil {
		fmt.Printf("an error occurred while fetching user: %v", err)
	}

	return model, err
}