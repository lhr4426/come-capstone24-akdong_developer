package routes

import (
	"capstone.com/controllers"
	"github.com/gin-gonic/gin"
)

func MapRoute(router *gin.Engine) {
	// map 관련 모든 routes 관리
	router.POST("/map_data", controllers.CreateMap())
	// router.GET("/map_data?map_id=:mapId&version=:mapVersion&chunk=:mapChunk", controllers.GetAMap())
	// router.GET("/map_data", controllers.GetAMap())
	// router.GET("/map_data/:mapId", controllers.GetAMap()) // map_id=&version=&chunck 변경 필요
}
