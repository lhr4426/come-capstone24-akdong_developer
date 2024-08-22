package routes

import (
	"github.com/gin-gonic/gin"
)

func GameRoute(router *gin.Engine) {
	router.GET("/creator_list", controllers.GetCreatorList())
}
