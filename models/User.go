package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

var Db *gorm.DB

type UserModel struct {
	ID        uint      `sql:"primary_key"`
	CreatedAt time.Time `sql:"default:CURRENT_TIMESTAMP"`
	UpdateAt  time.Time `sql:"default:CURRENT_TIMESTAMP"`
	EmailUser string    `sql:"size:200;unique_index"`
	GoogleId  string    `sql:"not null"`
	Picture   string    `sql:"null;type:text"`
}

func InsertFirstOnCreate(email string, google_id string, picture string) {
	user := UserModel{
		EmailUser: email,
		GoogleId:  google_id,
		Picture:   picture,
	}

	if err := Db.Where("email_user = ?", email).First(&user).Error; err != nil {
		// error handling...
		if gorm.IsRecordNotFoundError(err) {
			Db.Create(&user) // newUser not user
		}
	} else {
		Db.Model(&user).Where("email_user = ?", email).Update("picture", picture)
	}
}
