package hooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/services"
)

type BrowserNotificationsRegisterReq struct {
	RegistrationData string `json:"registration_data"`
}

func BrowserNotificationsRegister(
	browserNotifier *services.BrowserNotifier,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req BrowserNotificationsRegisterReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the subscriptions for the registration data
		err := browserNotifier.RegisterTarget(req.RegistrationData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the whoami info for this account
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{},
		})

	}
}
