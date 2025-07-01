package main

import (
	"github.com/Hritikpandey-ops/document-service/handlers"
	"github.com/Hritikpandey-ops/document-service/models"
	"github.com/Hritikpandey-ops/document-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	models.ConnectDatabase()

	utils.InitMinio()

	r := gin.Default()

	r.POST("/upload", handlers.UploadDocument)

	r.Run(":8083")
}

func init() {
	_ = godotenv.Load(".env")
}
