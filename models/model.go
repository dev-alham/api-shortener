package models

import (
	"github.com/jinzhu/gorm"
)

type ShortUrlModel struct {
	gorm.Model
	EmailUser	string `sql:"size:200;index;default:'anonymous'"`
	LongUrl		string `sql:"not null;size:200;unique_index"`
	ShortUrl 	string `sql:"not null;size:50;unique_index"`
	Status 		bool	`sql:"not null"`
}