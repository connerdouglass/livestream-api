package hooks

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/services"
)

func AppState(
	mainCreatorUsername string,
	telegramService *services.TelegramService,
	notificationsService *services.NotificationsService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the vapid keys for the notifications service
		var vapidPublicKey *string
		vapid, err := notificationsService.GetVapidKeyPair()
		if err != nil {
			fmt.Println("Error getting VAPID key: ", err.Error())
		}
		if vapid != nil {
			vapidPublicKey = &vapid.PublicKey
		}

		// Return the app state
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"main_creator_username": mainCreatorUsername,
				"telegram_bot_username": telegramService.BotUsername,
				"vapid_public_key":      vapidPublicKey,
			},
		})

	}
}
