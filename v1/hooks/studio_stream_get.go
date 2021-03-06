package hooks

import (
	"net/http"

	"github.com/connerdouglass/livestream-api/models"
	"github.com/connerdouglass/livestream-api/services"
	"github.com/connerdouglass/livestream-api/v1/utils"
	"github.com/gin-gonic/gin"
)

type StudioGetStreamReq struct {
	StreamID string `json:"stream_id"`
}

func StudioGetStream(
	creatorsService *services.CreatorsService,
	streamsService *services.StreamsService,
	membershipService *services.MembershipService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req StudioGetStreamReq
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

		// Return an empty response
		c.JSON(http.StatusOK, gin.H{
			"data": serializeStreamForStudio(stream),
		})

	}
}

func serializeStreamForStudio(stream *models.Stream) map[string]interface{} {
	if stream == nil {
		return nil
	}
	return map[string]interface{}{
		"id":                   stream.ID,
		"identifier":           stream.Identifier,
		"title":                stream.Title,
		"stream_key":           stream.StreamKey,
		"status":               stream.Status,
		"streaming":            stream.Streaming,
		"scheduled_start_date": stream.ScheduledStartDate.Unix(),
		"current_viewers":      stream.CurrentViewers,
	}
}
