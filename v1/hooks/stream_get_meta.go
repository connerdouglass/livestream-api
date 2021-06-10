package hooks

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/services"
)

type GetStreamMetaReq struct {
	StreamID string `json:"stream_id"`
}

func GetStreamMeta(
	streamsService *services.StreamsService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req GetStreamMetaReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Log the body
		fmt.Println("req: ", req)

	}
}
