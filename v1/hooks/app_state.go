package hooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/services"
)

func AppState(
	mainCreatorUsername string,
	telegramService *services.TelegramService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Return the app state
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"main_creator_username": mainCreatorUsername,
				"telegram_bot_username": telegramService.BotUsername,
			},
		})

	}
}
