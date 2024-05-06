package main

import (
	// "capstone.com/module/map/handler"

	"github.com/gin-gonic/gin"
	"capstone.com/configs"
	"capstone.com/routes"
	// "github.com/joho/godotenv"
	// "github.com/labstack/echo" // echo프레임워크 사용 // 높은 성능
)

func main() {

	router := gin.Default()

	configs.ConnectDB()

	routes.MapRoute(router)

	router.Run("localhost:8000")

	// router.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"data": "Hello from Gin-Gonic & mongoDB",
	// 	})
	// })

	// router.POST("/map", handler.Map)
	// router.Run("192.168.50.140:8000")

	// router.Run("localhost:8000")

}
