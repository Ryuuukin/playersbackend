// middlewares/adminAuth.go
package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAdmin(c *gin.Context) {
	isAdmin, exists := c.Get("isAdmin")
	if !exists || !isAdmin.(bool) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.Next()
}
