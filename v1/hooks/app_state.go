package hooks

import (
	"fmt"
	"net/http"

	"github.com/connerdouglass/livestream-api/services"
	"github.com/gin-gonic/gin"
)

func AppState(
	platformTitle string,
	mainCreatorUsername string,
	telegramService *services.TelegramService,
	browserNotifier *services.BrowserNotifier,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the vapid keys for the notifications service
		var vapidPublicKey *string
		vapid, err := browserNotifier.GetVapidKeyPair()
		if err != nil {
			fmt.Println("Error getting VAPID key: ", err.Error())
		}
		if vapid != nil {
			vapidPublicKey = &vapid.PublicKey
		}

		// Return the app state
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"platform_title":        platformTitle,
				"main_creator_username": mainCreatorUsername,
				"telegram_bot_username": telegramService.BotUsername,
				"vapid_public_key":      vapidPublicKey,
			},
		})

	}
}
