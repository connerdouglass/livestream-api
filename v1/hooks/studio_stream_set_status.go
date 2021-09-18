package hooks

import (
	"fmt"
	"net/http"
	"os"

	"github.com/connerdouglass/livestream-api/models"
	"github.com/connerdouglass/livestream-api/services"
	"github.com/connerdouglass/livestream-api/v1/utils"
	"github.com/gin-gonic/gin"
)

type StudioSetStreamStatusReq struct {
	StreamID string `json:"stream_id"`
	Status   string `json:"status"`
}

func StudioSetStreamStatus(
	streamsService *services.StreamsService,
	membershipService *services.MembershipService,
	notifier services.Notifier,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req StudioSetStreamStatusReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the account sending the request
		account := utils.CtxGetAccount(c)

		// Get the stream with the identifier
		stream, err := streamsService.GetStreamByIdentifier(req.StreamID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if stream == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "stream not found"})
			return
		}

		// Check if the account owns the stream
		isMember, err := membershipService.IsMember(stream.CreatorProfileID, account.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !isMember {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}

		// Update the status of the stream
		if err := streamsService.UpdateStatus(stream, req.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// If we're going live
		if req.Status == models.StreamStatus_Live {
			link := os.Getenv("TEMP_NOTIFY_LINK")
			var image *string
			if len(stream.CreatorProfile.Image) > 0 {
				image = &stream.CreatorProfile.Image
			}
			err := notifier.NotifySubscribers(
				stream.CreatorProfileID,
				&services.Notification{
					Title: stream.CreatorProfile.Name,
					Body:  fmt.Sprintf("%s just went live!", stream.CreatorProfile.Name),
					Link:  &link,
					Image: image,
				},
			)
			if err != nil {
				fmt.Println("Error sending notifications: ", err)
			}
		}

		// Return an empty response
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{},
		})

	}
}
