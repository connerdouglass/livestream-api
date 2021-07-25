package hooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/services"
)

type TelegramNotificationsUpdateSubReq struct {
	User       services.TelegramUser `json:"user"`
	CreatorID  uint64                `json:"creator_id"`
	Subscribed bool                  `json:"subscribed"`
}

func TelegramNotificationsUpdateSub(
	telegramNotifier *services.TelegramNotifier,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req TelegramNotificationsUpdateSubReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Subscribe to notifications
		if err := telegramNotifier.UpdateSub(
			&req.User,
			req.CreatorID,
			req.Subscribed,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the whoami info for this account
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{},
		})

	}
}
