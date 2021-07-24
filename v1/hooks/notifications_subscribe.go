package hooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/services"
)

type NotificationsSubscribeReq struct {
	RegistrationData string `json:"registration_data"`
	CreatorID        uint64 `json:"creator_id"`
}

func NotificationsSubscribe(
	notificationsService *services.NotificationsService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req NotificationsSubscribeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Subscribe to notifications
		if err := notificationsService.Subscribe(req.CreatorID, &req.RegistrationData, nil); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the whoami info for this account
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{},
		})

	}
}
