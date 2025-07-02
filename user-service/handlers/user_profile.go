package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/Hritikpandey-ops/user-service/database"
	"github.com/Hritikpandey-ops/user-service/models"
	"github.com/gin-gonic/gin"
)

var DB *sql.DB

// Create user profile
func CreateUserProfile(c *gin.Context) {
	var profile models.UserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	profile.CreatedAt = now
	profile.UpdatedAt = now

	query := `INSERT INTO user_profiles (user_id, full_name, phone, address, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := database.DB.QueryRow(query, profile.UserID, profile.FullName, profile.Phone, profile.Address, profile.CreatedAt, profile.UpdatedAt).Scan(&profile.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": profile})
}

// Get user profile by ID
func GetUserProfile(c *gin.Context) {
	id := c.Param("id")
	var profile models.UserProfile

	query := `SELECT id, user_id, full_name, phone, address, created_at, updated_at 
	          FROM user_profiles WHERE id = $1`
	err := database.DB.QueryRow(query, id).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FullName,
		&profile.Phone,
		&profile.Address,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": profile})
}

// Update user profile
func UpdateUserProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input models.UserProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	query := `UPDATE user_profiles SET full_name=$1, phone=$2, address=$3, updated_at=$4 WHERE id=$5`
	result, err := database.DB.Exec(query, input.FullName, input.Phone, input.Address, now, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}

	input.ID = uint(id)
	input.UpdatedAt = now
	c.JSON(http.StatusOK, gin.H{"data": input})
}

// Delete user profile
func DeleteUserProfile(c *gin.Context) {
	id := c.Param("id")

	query := `DELETE FROM user_profiles WHERE id = $1`
	result, err := database.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user profile"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User profile deleted successfully"})
}
