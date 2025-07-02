package main

import (
	"github.com/Hritikpandey-ops/user-service/database"
	"github.com/Hritikpandey-ops/user-service/handlers"
	middlewares "github.com/Hritikpandey-ops/user-service/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDatabase()
	r := gin.Default()

	protected := r.Group("/user")
	protected.Use(middlewares.AuthMiddleware())
	{
		protected.POST("/", handlers.CreateUserProfile)
		protected.GET("/:id", handlers.GetUserProfile)
		protected.PATCH("/:id", handlers.UpdateUserProfile)
		protected.DELETE("/:id", handlers.DeleteUserProfile)
	}

	r.Run(":8084")
}
