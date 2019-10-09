package controller

import (
	"api-shortener/config"
	"api-shortener/models"
	"api-shortener/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Msg struct{
	Status		bool	`json:"status"`
	Message		string	`json:"message"`
	Meta		Meta	`json:"meta"`
}

type Meta struct {
	LongUrl		string	`json:"long_url"`
	ShortUrl 	string	`json:"short_url"`
}

func CreateShortUrl(c *gin.Context) {
	unix_url := utils.RandStr(8)
	long_url := c.DefaultPostForm("long_url", "anoymous")
	tx := config.Db.Begin()

	shortUrl := models.ShortUrlModel{
		LongUrl:   long_url,
		ShortUrl:  unix_url,
	}
	config.Db.Where("long_url = ?", long_url).First(&shortUrl)
	if shortUrl.ID != 0{
		unix_url = shortUrl.ShortUrl
	}else{
		if err := tx.Save(&shortUrl).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, Msg{
				Status:  false,
				Message: "Please try again",
			})
			return
		}
		tx.Commit()
	}

	resp := Msg{
		true,
		"Create short url success",
		Meta{
			LongUrl:  long_url,
			ShortUrl: unix_url,
		},
	}

	c.JSON(http.StatusCreated, resp)
}
