package handlers

import (
	"net/http"

	"github.com/Hritikpandey-ops/vehicle-service/database"
	"github.com/Hritikpandey-ops/vehicle-service/models"
	"github.com/gin-gonic/gin"
)

func CreateVehicle(c *gin.Context) {
	var v models.Vehicle
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := database.DB.QueryRow(
		`INSERT INTO vehicles (brand, model, year, color, registration_number) 
         VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		v.Brand, v.Model, v.Year, v.Color, v.RegistrationNumber,
	).Scan(&v.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert vehicle"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vehicle created successfully", "vehicle_id": v.ID})
}

func GetVehicles(c *gin.Context) {
	rows, err := database.DB.Query(`SELECT id, brand, model, year, color, registration_number FROM vehicles`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vehicles"})
		return
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		var v models.Vehicle
		err := rows.Scan(&v.ID, &v.Brand, &v.Model, &v.Year, &v.Color, &v.RegistrationNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning vehicle data"})
			return
		}
		vehicles = append(vehicles, v)
	}

	c.JSON(http.StatusOK, vehicles)
}
