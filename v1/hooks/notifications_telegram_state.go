package hooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/services"
)

type TelegramNotificationsStateReq struct {
	User services.TelegramUser `json:"user"`
}

func TelegramNotificationsState(
	telegramNotifier *services.TelegramNotifier,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req TelegramNotificationsStateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the subscriptions for the registration data
		registered, subs, err := telegramNotifier.GetAllSubs(&req.User)
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
				"registered": registered,
				"subs":       subsSer,
			},
		})

	}
}
