package hooks

import (
	"net/http"

	"github.com/connerdouglass/livestream-api/services"
	"github.com/gin-gonic/gin"
)

type RtmpSetStreamingReq struct {
	StreamID  string `json:"stream_id"`
	Streaming bool   `json:"streaming"`
}

func RtmpSetStreaming(
	streamsService *services.StreamsService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req RtmpSetStreamingReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the details for the stream
		stream, err := streamsService.GetStreamByIdentifier(req.StreamID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if stream == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stream not found"})
			return
		}

		// Update the stream status
		if err := streamsService.UpdateStreaming(stream, req.Streaming); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Return a response of data for the stream
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"stream_id": stream.Identifier,
				"status":    stream.Status,
			},
		})

	}
}
