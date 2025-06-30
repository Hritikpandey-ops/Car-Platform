package middlewares

import (
	"net/http"
	"strings"

	"github.com/Hritikpandey-ops/auth-service/database"
	"github.com/Hritikpandey-ops/auth-service/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify JWT
		claims, err := utils.VerifyJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Fetch user role from DB using email
		var role string
		err = database.DB.QueryRow("SELECT role FROM users WHERE email = $1", claims.Email).Scan(&role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user role"})
			c.Abort()
			return
		}

		// Attach email and role to context
		c.Set("email", claims.Email)
		c.Set("role", role)

		c.Next()
	}
}
