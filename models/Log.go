package models

import (
	"api-shortener/config"
	"time"
)

type LogModel struct {
	ID        uint      `sql:"primary_key"`
	CreatedAt time.Time `sql:"default:CURRENT_TIMESTAMP"`
	ShortUrl  string    `sql:"not null; size:50"`
	IpAddr    string    `sql:"not null; index; type:char(40)"`
}

func InsertLog(log LogModel) (err error) {
	tx := config.Db.Begin()
	if tx.Error != nil {
		return err
	}

	if err := tx.Save(&log).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func GetCountRequestByDate(ip_addr string) int {
	var result int
	config.Db.Model(&LogModel{}).
		Where("DATE(created_at) = DATE(?)", time.Now()).
		Where("ip_addr = ?", ip_addr).
		Count(&result)
	return result
}
