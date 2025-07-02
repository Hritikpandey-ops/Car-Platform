package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Hritikpandey-ops/document-service/models"
	"github.com/gin-gonic/gin"
)

type Document struct {
	ID        int    `json:"id"`
	Filename  string `json:"filename"`
	VehicleID int    `json:"vehicle_id"`
}

func UpdateDocument(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var update struct {
		Filename  string `json:"filename"`
		VehicleID *int   `json:"vehicle_id"`
	}

	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if document exists
	var existingDoc Document
	err = models.DB.QueryRow("SELECT id, filename, vehicle_id FROM documents WHERE id = $1", id).
		Scan(&existingDoc.ID, &existingDoc.Filename, &existingDoc.VehicleID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch document"})
		return
	}

	// Update values conditionally
	if update.Filename != "" {
		existingDoc.Filename = update.Filename
	}
	if update.VehicleID != nil {
		existingDoc.VehicleID = *update.VehicleID
	}

	// Perform the update
	_, err = models.DB.Exec(
		"UPDATE documents SET filename = $1, vehicle_id = $2 WHERE id = $3",
		existingDoc.Filename, existingDoc.VehicleID, id,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document updated", "data": existingDoc})
}
