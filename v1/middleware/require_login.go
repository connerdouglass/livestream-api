package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/v1/utils"
)

// RequireLogin creates a middleware function to require authentication on a hook
func RequireLogin() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the account from the context
		account := utils.CtxGetAccount(c)
		if account == nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Authentication failed",
			})
			return
		}

		// Move to the next successfully
		c.Next()

	}
}
