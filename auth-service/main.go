package main

import (
	"os"

	"github.com/Hritikpandey-ops/auth-service/database"
	"github.com/Hritikpandey-ops/auth-service/handlers"
	middlewares "github.com/Hritikpandey-ops/auth-service/middleware"
	"github.com/Hritikpandey-ops/auth-service/utils"
	"github.com/gin-gonic/gin"
)

func main() {

	utils.InitLogger()
	utils.Log.Info("Starting Auth Service...")

	// Initialize the database
	if err := database.InitDB(); err != nil {
		panic("Failed to connect to DB: " + err.Error())
	}

	r := gin.Default()

	// Public routes
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)
	r.GET("/verify", handlers.VerifyEmail)

	// Protected route with JWT middleware
	api := r.Group("/api")
	api.Use(middlewares.AuthMiddleware())
	{
		api.POST("/admin/register", middlewares.AdminOnly(), handlers.AdminRegister)
		api.GET("/me", func(c *gin.Context) {
			email, _ := c.Get("email")
			c.JSON(200, gin.H{"message": "Authenticated", "email": email})
		})

		api.GET("/users/:id", handlers.GetUserByID)
		api.PATCH("/users/:id", handlers.UpdateUser)
		api.DELETE("/users/:id", handlers.DeleteUser)
		api.GET("/users/search", handlers.SearchUsers)
	}

	adminRoutes := api.Group("/admin")
	adminRoutes.Use(middlewares.AdminOnly())
	{
		adminRoutes.GET("/users", handlers.GetAllUsers)
		adminRoutes.PUT("/users/:id/promote", handlers.PromoteToAdmin)
	}

	// Start server
	r.Run(":" + os.Getenv("PORT"))
}
