package handlers

import (
	"net/http"
	"strconv"

	"github.com/Hritikpandey-ops/document-service/models"
	"github.com/gin-gonic/gin"
)

func GetDocumentsByVehicleID(c *gin.Context) {
	vehicleIDStr := c.Param("id")
	vehicleID, err := strconv.Atoi(vehicleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	rows, err := models.DB.Query("SELECT id, filename, url, vehicle_id FROM documents WHERE vehicle_id = $1", vehicleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve documents"})
		return
	}
	defer rows.Close()

	var documents []models.Document

	for rows.Next() {
		var doc models.Document
		if err := rows.Scan(&doc.ID, &doc.Filename, &doc.URL, &doc.VehicleID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan document"})
			return
		}
		documents = append(documents, doc)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading documents"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents})
}
