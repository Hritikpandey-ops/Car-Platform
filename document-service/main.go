package main

import (
	"github.com/Hritikpandey-ops/document-service/handlers"
	"github.com/Hritikpandey-ops/document-service/middleware"
	"github.com/Hritikpandey-ops/document-service/models"
	"github.com/Hritikpandey-ops/document-service/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	models.Connect()
	utils.InitMinio()

	r := gin.Default()
	api := r.Group("/documents")
	{
		api.Use(middleware.AuthMiddleware())

		api.POST("/upload", handlers.UploadDocument)
		api.GET("/vehicle/:id", handlers.GetDocumentsByVehicleID)
		api.PATCH("/:id", handlers.UpdateDocument)
		r.PATCH("/documents/:id/file", handlers.UpdateDocumentFile)
		api.DELETE("/:id", handlers.DeleteDocument)
	}

	r.Run(":8083")
}
