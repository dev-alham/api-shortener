package controller

import (
	"api-shortener/config"
	"api-shortener/models"
	"api-shortener/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CreateShortUrl(c *gin.Context) {
	unix_url := utils.RandStr(8)
	long_url := c.PostForm("long_url")
	var email string

	if long_url == "" || utils.CheckStrUrl(long_url) == false {
		c.JSON(http.StatusNotAcceptable, utils.ErrMsg{
			Status:  false,
			Message: "Not acceptable",
		})
		return
	}

	user := Auth(c)
	if user != nil {
		email = user.Email
	}

	tx := config.Db.Begin()
	shortUrl := models.ShortUrlModel{
		LongUrl:   long_url,
		ShortUrl:  unix_url,
		EmailUser: email,
	}
	config.Db.Where("long_url = ?", long_url).First(&shortUrl)
	if shortUrl.ID != 0 && shortUrl.EmailUser == "" && email == "" {
		unix_url = shortUrl.ShortUrl
	} else {
		if err := tx.Save(&shortUrl).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, utils.ErrMsg{
				Status:  false,
				Message: "Please try again",
			})
			return
		}
		tx.Commit()
	}

	resp := utils.SuccessMsg{
		true,
		"Create short url success " + email,
		utils.Meta{
			LongUrl:  long_url,
			ShortUrl: unix_url,
		},
	}
	c.JSON(http.StatusCreated, resp)
}

func GetOneShortUrl(c *gin.Context) {
	shortUrl := models.ShortUrlModel{}
	param_short_url := c.Param("param_short_url")
	if param_short_url == "" {
		c.JSON(http.StatusNotAcceptable, utils.ErrMsg{
			Status:  false,
			Message: "Not acceptable",
		})
		return
	}

	param_short_url = strings.ToUpper(param_short_url)
	config.Db.Where("short_url = ?", param_short_url).First(&shortUrl)

	if shortUrl.ID == 0 {
		c.JSON(http.StatusNotAcceptable, utils.ErrMsg{
			Status:  false,
			Message: "Data not found",
		})
		return
	}

	resp := utils.SuccessMsg{
		Status:  true,
		Message: "Data found",
		Meta: utils.Meta{
			LongUrl:  shortUrl.LongUrl,
			ShortUrl: shortUrl.ShortUrl,
		},
	}
	c.JSON(http.StatusOK, resp)
}
