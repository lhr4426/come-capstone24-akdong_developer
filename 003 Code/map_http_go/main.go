package main

import(
	"log"

	"capstone.com/module/map/handler"

	"github.com/joho/godotenv"
	"github.com/labstack/echo" // echo프레임워크 사용 // 높은 성능
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err := echo.New()
	e.POST("/map", handler.Map)

	e.Logger.Fatal(e.Start("192.168.50.140:8000"))

}