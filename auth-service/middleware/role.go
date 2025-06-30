package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole ensures that the user has the specified role
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Role not found in token"})
			return
		}

		role, ok := roleVal.(string)
		if !ok || role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		c.Next()
	}
}
