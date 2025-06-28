package main

import (
	"github.com/Hritikpandey-ops/vehicle-service/database"
	"github.com/Hritikpandey-ops/vehicle-service/handlers"
	middlewares "github.com/Hritikpandey-ops/vehicle-service/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	r := gin.Default()

	r.Use(gin.Logger(), gin.Recovery())

	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/vehicles", handlers.CreateVehicle)
		protected.GET("/vehicles", handlers.GetVehicles)
		protected.GET("/vehicles/:id", handlers.GetVehicleByID)
		protected.PUT("/vehicles/:id", handlers.UpdateVehicle)
		protected.DELETE("/vehicles/:id", handlers.DeleteVehicle)
		protected.GET("/vehicles/search", handlers.SearchVehicles)
	}

	r.Run(":8082")
}
