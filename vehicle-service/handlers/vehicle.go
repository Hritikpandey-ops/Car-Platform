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

func GetVehicleByID(c *gin.Context) {
	id := c.Param("id")

	var vehicle models.Vehicle
	err := database.DB.QueryRow(`SELECT id, brand, model, year, color, registration_number FROM vehicles WHERE id = $1`, id).
		Scan(&vehicle.ID, &vehicle.Brand, &vehicle.Model, &vehicle.Year, &vehicle.Color, &vehicle.RegistrationNumber)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}

func UpdateVehicle(c *gin.Context) {
	id := c.Param("id")
	var vehicle models.Vehicle

	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	_, err := database.DB.Exec(`
		UPDATE vehicles SET brand=$1, model=$2, year=$3, color=$4, registration_number=$5 WHERE id=$6`,
		vehicle.Brand, vehicle.Model, vehicle.Year, vehicle.Color, vehicle.RegistrationNumber, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vehicle updated"})
}

func DeleteVehicle(c *gin.Context) {
	id := c.Param("id")

	_, err := database.DB.Exec(`DELETE FROM vehicles WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vehicle"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vehicle deleted"})
}

func SearchVehicles(c *gin.Context) {
	query := c.Query("q")

	rows, err := database.DB.Query(`
		SELECT id, brand, model, year, color, registration_number
		FROM vehicles
		WHERE brand ILIKE '%' || $1 || '%' OR model ILIKE '%' || $1 || '%'`, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search vehicles"})
		return
	}
	defer rows.Close()

	var vehicles []models.Vehicle
	for rows.Next() {
		var v models.Vehicle
		rows.Scan(&v.ID, &v.Brand, &v.Model, &v.Year, &v.Color, &v.RegistrationNumber)
		vehicles = append(vehicles, v)
	}

	c.JSON(http.StatusOK, vehicles)
}
