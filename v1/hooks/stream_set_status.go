package hooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/services"
	"github.com/godocompany/livestream-api/v1/utils"
)

type SetStreamStatusReq struct {
	CreatorID uint64 `json:"creator_id"`
	StreamID  string `json:"stream_id"`
	Status    string `json:"status"`
}

func SetStreamStatus(
	accountsService *services.AccountsService,
	creatorsService *services.CreatorsService,
	streamsService *services.StreamsService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req SetStreamStatusReq
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
		owns, err := accountsService.DoesAccountOwnStream(account, stream)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !owns {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}

		// Update the status of the stream
		if err := streamsService.UpdateStatus(stream, req.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return an empty response
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{},
		})

	}
}
