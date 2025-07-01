package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Hritikpandey-ops/document-service/models"
	"github.com/Hritikpandey-ops/document-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func UploadDocument(c *gin.Context) {
	// Parse uploaded file
	header, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	file, err := header.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read file"})
		return
	}
	defer file.Close()

	// Get vehicle_id
	vehicleIDStr := c.PostForm("vehicle_id")
	if vehicleIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vehicle_id is required"})
		return
	}

	id, err := strconv.Atoi(vehicleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle_id"})
		return
	}

	// Generate object name
	objectName := fmt.Sprintf("%s_%s", uuid.New().String(), header.Filename)
	contentType := header.Header.Get("Content-Type")

	// Upload to MinIO
	_, err = utils.MinioClient.PutObject(context.Background(),
		"documents",
		objectName,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		log.Println("Upload error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to MinIO"})
		return
	}

	// Save metadata to DB
	document := models.Document{
		Filename:    header.Filename,
		URL:         objectName,
		ContentType: contentType,
		VehicleID:   id,
	}

	if err := models.DB.Create(&document).Error; err != nil {
		log.Println("Database insert error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded and saved",
		"object":  document,
	})
}
