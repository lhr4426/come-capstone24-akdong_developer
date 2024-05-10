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
		//Addr: "8080",
	}
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
	router.Run("0.0.0.0:8000")
}
