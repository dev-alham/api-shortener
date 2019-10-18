package controller

import (
	"api-shortener/models"
	"api-shortener/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CreateShortUrl(c *gin.Context) {
	unix_url := utils.RandStr(8)
	long_url := c.PostForm("long_url")
	tokenString := c.Request.Header.Get("Authorization")
	var email string = "null"

	if long_url == "" || utils.CheckStrUrl(long_url) == false {
		c.JSON(http.StatusNotAcceptable, utils.ErrMsg{
			Status:  false,
			Message: "Not acceptable",
		})
		return
	}

	user, _ := utils.GetSession(tokenString)
	if user != nil {
		email = user.Email
	}

	shortUrl := models.GetOne(models.ShortUrlModel{
		EmailUser: email,
		LongUrl:   long_url,
	})
	if shortUrl.ID != 0 {
		unix_url = shortUrl.ShortUrl
	} else {
		shortUrl = models.ShortUrlModel{
			EmailUser: email,
			LongUrl:   long_url,
			ShortUrl:  unix_url,
		}
		sts_insert := models.InsertUrl(shortUrl)
		if sts_insert != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrMsg{
				Status:  false,
				Message: "Please try again",
			})
			return
		}
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
	param_short_url := c.Param("param_short_url")
	if param_short_url == "" {
		c.JSON(http.StatusNotAcceptable, utils.ErrMsg{
			Status:  false,
			Message: "Not acceptable",
		})
		return
	}

	param_short_url = strings.ToUpper(param_short_url)
	shortUrl := models.GetOne(models.ShortUrlModel{
		ShortUrl: param_short_url,
	})

	if shortUrl.ID == 0 {
		c.JSON(http.StatusNotAcceptable, utils.ErrMsg{
			Status:  false,
			Message: "Data not found",
		})
		return
	}

	sts_update_count := models.UpdateShortUrl(models.ShortUrlModel{
		ShortUrl: param_short_url,
	}, models.ShortUrlModel{
		UpdateAt: utils.GetCurrentTime(),
		Count:    shortUrl.Count + 1,
	})

	if sts_update_count != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrMsg{
			Status:  false,
			Message: "Please try again",
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
