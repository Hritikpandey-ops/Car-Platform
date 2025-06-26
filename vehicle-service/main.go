package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/Hritikpandey-ops/vehicle-service/database"
	"github.com/Hritikpandey-ops/vehicle-service/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	database.Connect()

	router := gin.Default()
	router.POST("/vehicles", handlers.CreateVehicle)
	router.GET("/vehicles", handlers.GetVehicles)

	router.Run(":8082")
}
