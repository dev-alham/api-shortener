package config

import (
	"api-shortener/models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	DbHost = "localhost"
	DbPort = "3307"
	DbUser = "root"
	DbPass = ""
	DbName = "short_url"
)

var Db *gorm.DB
var err error

func InitConfig()  {
	dbCon := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		DbUser, DbPass, DbHost, DbPort, DbName)

	Db, err = gorm.Open("mysql", dbCon)
	if err != nil{
		panic("Failed connection database")
	}

	//Drops table if already exists
	//Db.DropTableIfExists(&models.ShortUrlModel{})

	//Auto create table based on Model
	//Db.AutoMigrate(&models.ShortUrlModel{})

	Db.AutoMigrate(&models.ShortUrlModel{})
}

func CloseConfig() error {
	return Db.Close()
}
