package hooks

import (
	"net/http"

	"github.com/connerdouglass/livestream-api/services"
	"github.com/connerdouglass/livestream-api/v1/utils"
	"github.com/gin-gonic/gin"
)

type StudioUpdateStreamReq struct {
	StreamID string                 `json:"stream_id"`
	Updates  services.StreamUpdates `json:"updates"`
}

func StudioUpdateStream(
	creatorsService *services.CreatorsService,
	streamsService *services.StreamsService,
	membershipService *services.MembershipService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req StudioUpdateStreamReq
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

		// Update the stream
		if err := streamsService.UpdateStream(stream, &req.Updates); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return an empty response
		c.JSON(http.StatusOK, gin.H{
			"data": serializeStreamForStudio(stream),
		})

	}
}
