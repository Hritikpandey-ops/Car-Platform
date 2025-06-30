package handlers

import (
	"net/http"

	"github.com/Hritikpandey-ops/auth-service/database"
	"github.com/Hritikpandey-ops/auth-service/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AdminRegister(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.Role != "admin" && input.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role must be 'user' or 'admin'"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	verificationToken := utils.GenerateVerificationToken()

	var userID int
	err := database.DB.QueryRow(`
		INSERT INTO users (email, password, verification_token, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		input.Email, string(hashedPassword), verificationToken, input.Role,
	).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	go utils.SendVerificationEmail(input.Email, verificationToken)

	c.JSON(http.StatusOK, gin.H{"message": "User created by admin", "user_id": userID})
}

func PromoteToAdmin(c *gin.Context) {
	id := c.Param("id")

	// Optionally, check if requester is an admin before allowing
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can promote users"})
		return
	}

	_, err := database.DB.Exec("UPDATE users SET role = 'admin' WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to promote user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin"})
}
