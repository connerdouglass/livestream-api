package hooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/services"
)

type BrowserNotificationsStateReq struct {
	RegistrationData string `json:"registration_data"`
}

func BrowserNotificationsState(
	browserNotifier *services.BrowserNotifier,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req BrowserNotificationsStateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the subscriptions for the registration data
		subs, err := browserNotifier.GetAllSubs(req.RegistrationData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Serialize all of the subs
		subsSer := make([]map[string]interface{}, len(subs))
		for i, sub := range subs {
			subsSer[i] = map[string]interface{}{
				"creator_id": sub.CreatorProfileID,
			}
		}

		// Return the whoami info for this account
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"registered": true,
				"subs":       subsSer,
			},
		})

	}
}
