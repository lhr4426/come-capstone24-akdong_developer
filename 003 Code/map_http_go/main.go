package main

import (
	"capstone.com/configs"
	"capstone.com/routes"
	// "capstone.com/controllers"
	"github.com/gin-gonic/gin"
	// "github.com/labstack/echo" // echo프레임워크 사용 // 높은 성능
)

func main() {

	router := gin.Default()

	configs.ConnectDB()

	routes.MapRoute(router)

	router.Run("localhost:8000")

	// router.Run("192.168.50.140:8000")


}
