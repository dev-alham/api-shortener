package main

import (
	"api-shortener/config"
	"api-shortener/controller"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	router := gin.Default()
	os.Setenv("TZ", "Asia/Jakarta")
	config.InitConfig()

	v1 := router.Group("/api/v1/")
	{
		v1.POST("/", controller.CreateShortUrl)
		v1.GET("/:param_short_url", controller.GetOneShortUrl)
	}

	router.Run(":3000")

}
