package middleware

import (
	"net/http"

	"github.com/connerdouglass/livestream-api/services"
	"github.com/gin-gonic/gin"
)

// RequireRtmpAuth creates a middleware function to require RTMP auth on a hook
func RequireRtmpAuth(
	rtmpAuthService *services.RtmpAuthService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Check for the bearer token
		token := c.GetString("bearer_token")

		// Validate the token
		if !rtmpAuthService.CheckPasscode(token) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Authentication failed",
			})
			return
		}

		// Move to the next successfully
		c.Next()

	}
}
