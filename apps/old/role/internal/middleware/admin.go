package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	permissionapp "mandacode.com/accounts/role/internal/app/permission"
)

func NewAdminMiddleware(permissionApp permissionapp.PermissionApp, uidHeader string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user ID from the request header
		userID := c.GetHeader(uidHeader)
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: user ID not provided"})
			return
		}
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
			return
		}

		// Check if the user is an admin
		isAdmin, err := permissionApp.CheckAdmin(userUUID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}
		if !isAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: user is not an admin"})
			return
		}

		// If the user is an admin, proceed with the request
		c.Next()
	}
}
