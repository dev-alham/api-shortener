package controller

import (
	"api-shortener/cache"
	"api-shortener/models"
	"api-shortener/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var CACHE_DIR_SHORT_URL = "SHORT URL"
var CACHE_DIR_LONG_URL = "LONG URL"
var CACHE_DIR_LIMIT = "LIMIT"
var LIMIT_REQUEST_GET_DAY = 10
var LIMIT_REQUEST_POST_DAY = 3

func CreateShortUrl(c *gin.Context) {
	unix_url := utils.RandStr(8)
	long_url_req := c.PostForm("long_url")
	short_url_req := c.PostForm("short_url")
	tokenString := c.Request.Header.Get("Authorization")
	var email string = "null"
	shortUrl := models.ShortUrlModel{}

	short_url_req = strings.ToUpper(short_url_req)
	user, _ := utils.GetSession(tokenString)
	long_url_req = utils.DeletePrefixUrl(long_url_req)
	ip_addr := c.ClientIP()

	// validate long url
	if long_url_req == "" || utils.CheckStrUrl(long_url_req) == false {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.ErrMsg{
			Status:  false,
			Message: "Not acceptable",
		})
		return
	}

	// validate custom url not email
	if user == nil && short_url_req != "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.ErrMsg{
			Status:  false,
			Message: "Not acceptable custom url",
		})
		return
	}

	// validate custom short url
	if utils.ValidateBetween(len(short_url_req), 4, 8) == false &&
		short_url_req != "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.ErrMsg{
			Status:  false,
			Message: "Length short url not available",
		})
		return
	}

	// validate email
	if user != nil {
		// change anonymous to email user
		email = user.Email
		if short_url_req != "" {
			// change unix url to custom short url
			unix_url = short_url_req
		}
	}

	if user == nil {
		sts_limit, _ := cache.GetValue(CACHE_DIR_LIMIT+":"+CACHE_DIR_SHORT_URL, ip_addr)
		if sts_limit == nil {
			count_req_anonym := models.GetCountAnonymousRequest(ip_addr, email)
			if count_req_anonym >= LIMIT_REQUEST_POST_DAY {
				cache.SetValueWithTTL(CACHE_DIR_LIMIT+":"+CACHE_DIR_SHORT_URL, ip_addr, true, 60)
				sts_limit = true
			}
		}

		if sts_limit == true {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, utils.ErrMsg{
				Status:  false,
				Message: "Limit request",
			})
			return
		}
	}

	// get short url in redis
	cache_short_url, _ := cache.GetValue(CACHE_DIR_SHORT_URL+":"+email, long_url_req)
	if cache_short_url == nil {
		shortUrl = models.GetOne(models.ShortUrlModel{
			EmailUser: email,
			LongUrl:   long_url_req,
		})
		if shortUrl.ID != 0 {
			cache.SetValueWithTTL(CACHE_DIR_SHORT_URL+":"+email, long_url_req, shortUrl.ShortUrl, 60)
			cache_short_url = shortUrl.ShortUrl
		} else {
			cache.SetValueWithTTL(CACHE_DIR_SHORT_URL+":"+email, long_url_req, unix_url, 5)
		}
	}

	// check user not used email
	if cache_short_url != nil {
		c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMsg{
			Status:  true,
			Message: "Shortl url success",
			Meta: utils.Meta{
				LongUrl:  long_url_req,
				ShortUrl: fmt.Sprintf("%v", cache_short_url),
			},
		})
		return
	}

	shortUrl = models.GetOne(models.ShortUrlModel{
		EmailUser: email,
		ShortUrl:  unix_url,
	})

	if shortUrl.ID != 0 {
		// check email and short url
		if shortUrl.EmailUser == user.Email &&
			shortUrl.ShortUrl == unix_url {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.ErrMsg{
				Status:  false,
				Message: "Custom short url already used",
			})
			cache.DelKey(CACHE_DIR_SHORT_URL+":"+email, long_url_req)
			return
		}
	}

	sts_insert := models.InsertUrl(models.ShortUrlModel{
		EmailUser: email,
		LongUrl:   long_url_req,
		ShortUrl:  unix_url,
		IpAddr:    ip_addr,
	})
	if sts_insert != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrMsg{
			Status:  false,
			Message: "Please try again",
		})
		cache.DelKey(CACHE_DIR_SHORT_URL+":"+email, long_url_req)
		return
	}

	c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.SuccessMsg{
		Status:  true,
		Message: "Short url success",
		Meta: utils.Meta{
			LongUrl:  long_url_req,
			ShortUrl: unix_url,
		},
	})
}

func GetOneShortUrl(c *gin.Context) {
	short_url_req := c.Param("param_short_url")

	ip_addr := c.ClientIP()
	short_url_req = strings.ToUpper(short_url_req)
	shortUrl := models.ShortUrlModel{}

	if short_url_req == "" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.ErrMsg{
			Status:  false,
			Message: "Not acceptable",
		})
		return
	}

	sts_limit, _ := cache.GetValue(CACHE_DIR_LIMIT+":"+CACHE_DIR_LONG_URL, ip_addr)
	if sts_limit == nil {
		cache_limit_req := models.GetCountRequestByDate(ip_addr)
		if cache_limit_req >= LIMIT_REQUEST_GET_DAY {
			cache.SetValueWithTTL(CACHE_DIR_LIMIT+":"+CACHE_DIR_LONG_URL, ip_addr, true, 60)
			sts_limit = true
		}
	}

	if sts_limit == true {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, utils.ErrMsg{
			Status:  false,
			Message: "Limit request",
		})
		return
	}

	insert_log := models.InsertLog(models.LogModel{
		ShortUrl: short_url_req,
		IpAddr:   ip_addr,
	})

	if insert_log != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrMsg{
			Status:  false,
			Message: "Please try again",
		})
	}

	cache_long_url, _ := cache.GetValue(CACHE_DIR_LONG_URL, short_url_req)
	if cache_long_url == nil {
		shortUrl = models.GetOne(models.ShortUrlModel{
			ShortUrl: short_url_req,
		})

		if shortUrl.ID != 0 {
			cache_long_url = shortUrl.LongUrl
		} else {
			cache_long_url = "-"
		}
		cache.SetValueWithTTL(CACHE_DIR_LONG_URL, short_url_req, cache_long_url, 60)
	}

	if fmt.Sprintf("%v", cache_long_url) == "-" {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, utils.ErrMsg{
			Status:  false,
			Message: "Data not found",
		})
		return
	}

	sts_update_count := models.UpdateShortUrl(models.ShortUrlModel{
		ShortUrl: short_url_req,
	}, models.ShortUrlModel{
		UpdateAt: utils.GetCurrentTime(),
		Count:    shortUrl.Count + 1,
	})
	if sts_update_count != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ErrMsg{
			Status:  false,
			Message: "Please try again",
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, utils.SuccessMsg{
		Status:  true,
		Message: "Data found",
		Meta: utils.Meta{
			LongUrl:  fmt.Sprintf("%v", cache_long_url),
			ShortUrl: short_url_req,
		},
	})
}
