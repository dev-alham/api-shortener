package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

var Db *gorm.DB
var err error

func DbInit() {
	dbCon := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	Db, err = gorm.Open("mysql", dbCon)
	if err != nil {
		panic("Failed connection database")
	}

	// print query
	Db.LogMode(PRINT_QUERY)

	//Drops table if already exists
	//Db.DropTableIfExists(&models.ShortUrlModel{})

	//Auto create table based on Model
	//Db.AutoMigrate(&models.ShortUrlModel{})

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return PREFIX_DB_NAME + "_" + defaultTableName
	}
}

func CloseConfig() error {
	return Db.Close()
}
