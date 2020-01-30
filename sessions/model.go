package sessions

import (
	"time"
	"workour-api/common"
)

type Session struct {
	ID		uint 		`gorm:"primary_key"`
	Email	string 		`gorm:"column:email;type:varchar(100)"`
	Expiry	time.Time	`gorm:"column:expiry;type:timestamptz"`
	UserId	int			`gorm:"type:int REFERENCES users(id)"`
}

func (s Session) Create(data interface{}) (uint, error) {
	db := common.GetDB()
	session := data.(Session)

	err := db.Create(&session).Error

	if err != nil {
		return 0, err
	}

	return session.ID, nil
}

func (s Session) GetById(id uint) (Session, error) {
	db := common.GetDB()
	session := Session{}

	err := db.Where("id = ?", id).Find(&session).Error

	if err != nil {
		return Session{}, err
	}

	return session, nil
}

func (s Session) Update(data interface{}) error {
	db := common.GetDB()
	session := data.(Session)

	err := db.Save(&session).Error

	if err != nil {
		return err
	}

	return nil
}

func (s Session) Delete(id uint) error {
	db := common.GetDB()
	session := Session{ID: id}

	err := db.Unscoped().Delete(&session).Error

	if err != nil {
		return err
	}

	return nil
}