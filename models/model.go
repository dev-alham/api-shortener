package models

import (
	"time"
)

type ShortUrlModel struct {
	ID        uint      `sql:"primary_key"`
	CreatedAt time.Time `sql:"default:CURRENT_TIMESTAMP"`
	EmailUser string    `sql:"size:200;index"`
	LongUrl   string    `sql:"not null;size:200"`
	ShortUrl  string    `sql:"not null;size:50;unique_index"`
}
