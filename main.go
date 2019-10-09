package main

import (
	"api-shortener/config"
	"api-shortener/controller"
	"github.com/gin-gonic/gin"
)

func main()  {
	router := gin.Default()
	config.InitConfig()

	v1 := router.Group("/api/v1/")
	{
		v1.POST("/", controller.CreateShortUrl)
	}

	router.Run(":3000")


}
