package models

import (
	"api-shortener/config"
	"time"
)

type ShortUrlModel struct {
	ID        uint      `sql:"primary_key"`
	CreatedAt time.Time `sql:"default:CURRENT_TIMESTAMP"`
	EmailUser string    `sql:"size:200;index"`
	LongUrl   string    `sql:"not null;size:200"`
	ShortUrl  string    `sql:"not null;size:50;unique_index"`
}

func GetOne(column string, value string) (shortUrl ShortUrlModel) {
	config.Db.Where(column+"= ?", value).First(&shortUrl)
	return
}

func MultipleCondition(i interface{}) (shortUrl ShortUrlModel) {
	config.Db.Where(i).Find(&shortUrl)
	return
}

func InsertUrl(shortUrl ShortUrlModel) error {
	tx := config.Db.Begin()
	if err := tx.Save(&shortUrl).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
