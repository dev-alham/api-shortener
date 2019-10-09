package controller

import (
	"api-shortener/config"
	"api-shortener/models"
	"api-shortener/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Msg struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Meta    Meta   `json:"meta"`
}

type ErrMsg struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type Meta struct {
	LongUrl  string `json:"long_url"`
	ShortUrl string `json:"short_url"`
}

func CreateShortUrl(c *gin.Context) {
	unix_url := utils.RandStr(8)
	long_url := c.PostForm("long_url")

	if long_url == "" || utils.CheckStrUrl(long_url) == false {
		c.JSON(http.StatusNotAcceptable, ErrMsg{
			Status:  false,
			Message: "Not acceptable",
		})
		return
	}

	tx := config.Db.Begin()
	shortUrl := models.ShortUrlModel{
		LongUrl:  long_url,
		ShortUrl: unix_url,
	}
	config.Db.Where("long_url = ?", long_url).First(&shortUrl)
	if shortUrl.ID != 0 && shortUrl.EmailUser == "" {
		unix_url = shortUrl.ShortUrl
	} else {
		if err := tx.Save(&shortUrl).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, ErrMsg{
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

func GetOneShortUrl(c *gin.Context) {
	shortUrl := models.ShortUrlModel{}
	param_short_url := c.Param("param_short_url")
	if param_short_url == "" {
		c.JSON(http.StatusNotAcceptable, ErrMsg{
			Status:  false,
			Message: "Not acceptable",
		})
		return
	}
	config.Db.Where("short_url = ?", param_short_url).First(&shortUrl)

	if shortUrl.ID == 0 {
		c.JSON(http.StatusNotAcceptable, ErrMsg{
			Status:  false,
			Message: "Data not found",
		})
		return
	}

	resp := Msg{
		Status:  true,
		Message: "Data found",
		Meta: Meta{
			LongUrl:  shortUrl.LongUrl,
			ShortUrl: shortUrl.ShortUrl,
		},
	}

	c.JSON(http.StatusOK, resp)
}
