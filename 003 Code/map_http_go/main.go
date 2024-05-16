package main

import (
	"map_http_go/configs"
	"map_http_go/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	configs.ConnectDB()

	routes.MapRoute(router)

	// router.Run("localhost:8000")

	router.Run("0.0.0.0:8080")

}
