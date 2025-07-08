package main

import (
	"log"
	"os"

	"github.com/Hritikpandey-ops/mail-service/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}

	r := gin.Default()
	r.POST("/send", handlers.SendEmailHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	r.Run(":" + port)
}
