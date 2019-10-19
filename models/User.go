package models

import (
	"api-shortener/config"
	"github.com/jinzhu/gorm"
	"time"
)

var Db *gorm.DB

type UserModel struct {
	ID        uint      `sql:"primary_key"`
	CreatedAt time.Time `sql:"default: CURRENT_TIMESTAMP"`
	UpdateAt  time.Time `sql:"default: CURRENT_TIMESTAMP"`
	EmailUser string    `sql:"size:200; unique_index"`
	GoogleId  string    `sql:"not null; size:200;"`
	Picture   string    `sql:"null; type:text"`
}

func InsertFirstOnCreate(where_user UserModel, update_user UserModel) (err error) {
	tx := config.Db.Begin()

	if tx.Error != nil {
		return err
	}

	if err := tx.Model(&UserModel{}).
		Where(where_user).
		Assign(&update_user).
		FirstOrCreate(&UserModel{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
