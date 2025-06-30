package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Hritikpandey-ops/auth-service/models"
	"github.com/Hritikpandey-ops/auth-service/utils"

	"github.com/Hritikpandey-ops/auth-service/database"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash the password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	verificationToken := utils.GenerateVerificationToken()

	role := "user"
	err := database.DB.QueryRow(`
		INSERT INTO users (email, password, verification_token, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		user.Email, string(hashedPassword), verificationToken, role,
	).Scan(&user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Send email
	go utils.SendVerificationEmail(user.Email, verificationToken)

	c.JSON(http.StatusOK, gin.H{
		"message": "User created, please verify your email",
		"user_id": user.ID,
	})
}

func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token"})
		return
	}

	result, err := database.DB.Exec(`
		UPDATE users
		SET is_verified = true, verification_token = NULL
		WHERE verification_token = $1 AND is_verified = false`, token)

	rowsAffected, _ := result.RowsAffected()

	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully "})
}

func Login(c *gin.Context) {
	var input models.User
	var dbUser models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := database.DB.QueryRow(
		"SELECT id, email, password, is_verified, role FROM users WHERE email=$1",
		input.Email,
	).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Password, &dbUser.IsVerified, &dbUser.Role)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	if !dbUser.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not verified"})
		return
	}

	token, err := utils.GenerateJWT(dbUser.Email, dbUser.Role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Get all users
func GetAllUsers(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(403, gin.H{"error": "Only admins can access this resource"})
		return
	}

	rows, err := database.DB.Query(`SELECT id, email, is_verified, role FROM users`)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.IsVerified, &user.Role); err != nil {
			utils.Log.WithError(err).Error("Failed to Get user row")
			continue
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": "Error reading users from database"})
		return
	}

	c.JSON(200, users)
}

// Get user by ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	err := database.DB.QueryRow(`SELECT id, email, is_verified, role FROM users WHERE id=$1`, id).Scan(
		&user.ID, &user.Email, &user.IsVerified, &user.Role)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, user)
}

// Update user
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Email      *string `json:"email"`
		IsVerified *bool   `json:"is_verified"`
		Role       *string `json:"role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	// Validate role if provided
	if input.Role != nil {
		validRoles := map[string]bool{"user": true, "admin": true}
		if !validRoles[*input.Role] {
			c.JSON(400, gin.H{"error": "Invalid role"})
			return
		}
	}

	// Build the SQL dynamically
	query := "UPDATE users SET"
	args := []interface{}{}
	i := 1

	if input.Email != nil {
		query += fmt.Sprintf(" email=$%d,", i)
		args = append(args, *input.Email)
		i++
	}
	if input.IsVerified != nil {
		query += fmt.Sprintf(" is_verified=$%d,", i)
		args = append(args, *input.IsVerified)
		i++
	}
	if input.Role != nil {
		query += fmt.Sprintf(" role=$%d,", i)
		args = append(args, *input.Role)
		i++
	}

	// Remove trailing comma
	if len(args) == 0 {
		c.JSON(400, gin.H{"error": "No fields to update"})
		return
	}
	query = query[:len(query)-1]

	// Add WHERE clause
	query += fmt.Sprintf(" WHERE id=$%d", i)
	args = append(args, id)

	_, err := database.DB.Exec(query, args...)
	if err != nil {
		utils.Log.WithError(err).Error("Failed to update user")
		c.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(200, gin.H{"message": "User updated"})
}

// Delete user
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	_, err := database.DB.Exec(`DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted"})
}

// Search users (by email)
func SearchUsers(c *gin.Context) {
	query := c.Query("q")
	rows, err := database.DB.Query("SELECT id, email, is_verified, role FROM users WHERE email ILIKE '%' || $1 || '%'", query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.IsVerified, &user.Role); err != nil {
			continue
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}
