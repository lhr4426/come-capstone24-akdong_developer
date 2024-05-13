package main

import (

	// "module/bin"

	"log"
	"net/http"

	"capstone.com/module/handler"
	"github.com/gin-gonic/gin"
)

func main() {

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

	router.Run("192.168.50.140:8000")
}
