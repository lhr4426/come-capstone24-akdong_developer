package main

import (
	"asset_http_go/configs"
	"asset_http_go/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 데이터베이스 실행
	configs.ConnectDB()

	// 라우트
	routes.AssetRoute(router)

	router.Run("0.0.0.0:8070")
}
