package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-KEY")
		secret := os.Getenv("INTERNAL_API_KEY")

		if key != secret {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			c.Abort() // Stop the request here
			return
		}
		c.Next()
	}
}
