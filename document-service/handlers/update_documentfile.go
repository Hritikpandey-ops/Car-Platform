package handlers

import (
	"context"
	"database/sql"
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

func UpdateDocumentFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	// Retrieve existing document
	var doc models.Document
	err = models.DB.QueryRow("SELECT id, filename, url, content_type, vehicle_id FROM documents WHERE id = $1", id).
		Scan(&doc.ID, &doc.Filename, &doc.URL, &doc.ContentType, &doc.VehicleID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch document"})
		return
	}

	// Get new file from request
	header, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New file is required"})
		return
	}

	file, err := header.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read new file"})
		return
	}
	defer file.Close()

	// Delete old file from MinIO
	err = utils.MinioClient.RemoveObject(context.Background(), "documents", doc.URL, minio.RemoveObjectOptions{})
	if err != nil {
		log.Println("Warning: failed to delete old object from MinIO:", err)
	}

	// Upload new file
	newObjectName := fmt.Sprintf("%s_%s", uuid.New().String(), header.Filename)
	contentType := header.Header.Get("Content-Type")

	_, err = utils.MinioClient.PutObject(
		context.Background(),
		"documents",
		newObjectName,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload new file to MinIO"})
		return
	}

	// Update DB record
	_, err = models.DB.Exec(
		"UPDATE documents SET filename = $1, url = $2, content_type = $3 WHERE id = $4",
		header.Filename, newObjectName, contentType, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document file updated successfully",
		"data": gin.H{
			"id":           id,
			"filename":     header.Filename,
			"url":          newObjectName,
			"content_type": contentType,
		},
	})
}
