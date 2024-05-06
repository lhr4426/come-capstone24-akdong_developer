package routes

import (
	"asset_http_go/controllers"

	"github.com/gin-gonic/gin"
)

func AssetRoute(router *gin.Engine) {
	router.POST("/asset_upload", controllers.CreateAsset())
	router.GET("/asset_info", controllers.GetAsset())
}
