package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Hritikpandey-ops/document-service/models"
	"github.com/Hritikpandey-ops/document-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)


func DeleteDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var doc models.Document
	err = models.DB.QueryRow("SELECT id, url, filename FROM documents WHERE id = $1", id).
		Scan(&doc.ID, &doc.URL, &doc.Filename)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch document"})
		return
	}

	// Remove from MinIO
	err = utils.MinioClient.RemoveObject(context.Background(), "documents", doc.URL, minio.RemoveObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete from MinIO"})
		return
	}

	// Delete from DB
	_, err = models.DB.Exec("DELETE FROM documents WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete document from DB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
