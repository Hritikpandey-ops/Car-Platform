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
