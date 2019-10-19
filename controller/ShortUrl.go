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
	long_url_req := c.PostForm("long_url")
	short_url_req := c.PostForm("short_url")
	tokenString := c.Request.Header.Get("Authorization")
	var email string = "null"

	short_url_req = strings.ToUpper(short_url_req)
	user, _ := utils.GetSession(tokenString)

	if long_url_req == "" || utils.CheckStrUrl(long_url_req) == false {
		c.JSON(http.StatusNotAcceptable, utils.ErrMsg{
			Status:  false,
			Message: "Not acceptable",
		})
		return
	}

	if user != nil {
		email = user.Email
		if short_url_req != "" {
			unix_url = short_url_req

			// validate custom short url
			if utils.ValidateBetween(len(short_url_req), 4, 8) == false {
				c.JSON(http.StatusNotAcceptable, utils.ErrMsg{
					Status:  false,
					Message: "Length short url not available",
				})
				return
			}

			check_short_url := models.GetOne(models.ShortUrlModel{
				EmailUser: email,
				ShortUrl:  unix_url,
			})

			if check_short_url.EmailUser == user.Email &&
				check_short_url.ShortUrl == unix_url &&
				check_short_url.LongUrl == long_url_req {
				// check all by email, short url and long url
				c.JSON(http.StatusNotAcceptable, utils.SuccessMsg{
					Status:  true,
					Message: "Create short url success " + email,
					Meta: utils.Meta{
						LongUrl:  long_url_req,
						ShortUrl: unix_url,
					},
				})
				return
			} else if check_short_url.EmailUser == user.Email &&
				check_short_url.ShortUrl == unix_url {
				// check email and short url
				c.JSON(http.StatusNotAcceptable, utils.ErrMsg{
					Status:  false,
					Message: "Short url already used",
				})
				return
			}

			// update custom url
			sts_update := models.UpdateShortUrl(models.ShortUrlModel{
				EmailUser: user.Email,
				LongUrl:   long_url_req,
			}, models.ShortUrlModel{
				ShortUrl: short_url_req,
			})

			if sts_update != nil {
				// check error
				c.JSON(http.StatusInternalServerError, utils.ErrMsg{
					Status:  false,
					Message: "Please try again",
				})
				return
			} else {
				// response update url
				c.JSON(http.StatusCreated, utils.SuccessMsg{
					Status:  true,
					Message: "Update short url",
					Meta: utils.Meta{
						LongUrl:  long_url_req,
						ShortUrl: short_url_req,
					},
				})
				return
			}
		}
	}

	// check user not used email
	shortUrl := models.GetOne(models.ShortUrlModel{
		EmailUser: email,
		LongUrl:   long_url_req,
	})

	if shortUrl.ID != 0 {
		unix_url = shortUrl.ShortUrl
	} else {
		shortUrl = models.ShortUrlModel{
			EmailUser: email,
			LongUrl:   long_url_req,
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
			LongUrl:  long_url_req,
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
