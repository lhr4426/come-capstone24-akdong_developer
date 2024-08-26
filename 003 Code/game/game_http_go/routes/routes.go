package routes

import (
	"game_http_go/controllers"

	"github.com/gin-gonic/gin"
)

func GameRoute(router *gin.Engine) {
	router.GET("/", controllers.GetTest())
	// router.GET("/creator_list", controllers.GetCreatorList())
}
