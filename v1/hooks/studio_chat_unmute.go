package hooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/services"
)

type StudioChatUnmuteReq struct {
	// CreatorID string `json:"creator_id"`
	Username string `json:"username"`
}

func StudioChatUnmute(
	accountsService *services.AccountsService,
	chatService *services.ChatService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req StudioChatUnmuteReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the account sending the request
		// account := utils.CtxGetAccount(c)

		// Mute the user on the chat
		if err := chatService.UnmuteUser(req.Username); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Otherwise return something successfully
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{},
		})

	}
}
