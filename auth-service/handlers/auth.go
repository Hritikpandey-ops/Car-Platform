package handlers

import (
	"database/sql"
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

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	verificationToken := utils.GenerateVerificationToken()

	err := database.DB.QueryRow(`
		INSERT INTO users (email, password, verification_token)
		VALUES ($1, $2, $3)
		RETURNING id`,
		user.Email, string(hashedPassword), verificationToken,
	).Scan(&user.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	// Send email
	go utils.SendVerificationEmail(user.Email, verificationToken)

	c.JSON(http.StatusOK, gin.H{"message": "User created, please verify email", "user_id": user.ID})
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

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully ✅"})
}

func Login(c *gin.Context) {
	var input models.User
	var dbUser models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := database.DB.QueryRow("SELECT id, password, is_verified FROM users WHERE email=$1", input.Email).
		Scan(&dbUser.ID, &dbUser.Password, &dbUser.IsVerified)

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

	token, err := utils.GenerateJWT(input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Get all users
func GetAllUsers(c *gin.Context) {
	rows, err := database.DB.Query(`SELECT id, email, is_verified FROM users`)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.IsVerified); err != nil {
			continue
		}
		users = append(users, user)
	}
	c.JSON(200, users)
}

// Get user by ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	err := database.DB.QueryRow(`SELECT id, email, is_verified FROM users WHERE id=$1`, id).Scan(
		&user.ID, &user.Email, &user.IsVerified)
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
		Email      string `json:"email"`
		IsVerified bool   `json:"is_verified"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	_, err := database.DB.Exec(`UPDATE users SET email=$1, is_verified=$2 WHERE id=$3`, input.Email, input.IsVerified, id)
	if err != nil {
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
	rows, err := database.DB.Query("SELECT id, email, is_verified FROM users WHERE email ILIKE '%' || $1 || '%'", query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search users"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Email, &user.IsVerified); err != nil {
			continue
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}
