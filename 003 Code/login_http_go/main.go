package main

import (

	// "module/bin"

	"log"
	"net/http"

	"capstone.com/module/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()
	router.POST("/signup", handler.SignUp())
	router.POST("/login", handler.LogIn())

	log.Println("start login server")
	server := &http.Server{
		// Addr: "8000",
	}
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}

	router.Run("0.0.0.0:8000")
}
