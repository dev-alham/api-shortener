package models

import (
	"api-shortener/config"
	"time"
)

type ShortUrlModel struct {
	ID        uint      `sql:"primary_key"`
	CreatedAt time.Time `sql:"default:CURRENT_TIMESTAMP"`
	UpdateAt  time.Time `sql:"default:CURRENT_TIMESTAMP"`
	EmailUser string    `sql:"size:200;default: null;index"`
	LongUrl   string    `sql:"not null;size:200"`
	ShortUrl  string    `sql:"not null;size:50;unique_index"`
	Count     int64     `sql:"not null; default: 0"`
}

func GetOne(short_url ShortUrlModel) (shortUrl ShortUrlModel) {
	/* example

	var email string = "null"

	shortUrl := models.GetOne(models.ShortUrlModel{
		EmailUser: email,
		LongUrl:   long_url,
	})
	*/

	config.Db.Where(&short_url).First(&shortUrl)
	return
}

func MultipleConditionAll(i interface{}) (shortUrl ShortUrlModel) {
	/* example

	multi_conditions := map[string]interface{}{
		"long_url":   long_url,
		"email_user": email,
	}
	*/

	config.Db.Where(i).Find(&shortUrl)
	return
}

func InsertUrl(shortUrl ShortUrlModel) (err error) {
	/* example

	sts_insert := models.InsertUrl(shortUrl)
	if sts_insert != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrMsg{
			Status:  false,
			Message: "Please try again",
		})
		return
	}
	*/
	tx := config.Db.Begin()
	if tx.Error != nil {
		return err
	}

	if err := tx.Save(&shortUrl).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func UpdateShortUrl(where_param ShortUrlModel, update_param ShortUrlModel) (err error) {
	tx := config.Db.Begin()
	if tx.Error != nil {
		return err
	}

	if err := tx.Model(&ShortUrlModel{}).
		Where(where_param).
		Update(&update_param).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
