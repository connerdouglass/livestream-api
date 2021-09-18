package hooks

import (
	"net/http"

	"github.com/connerdouglass/livestream-api/services"
	"github.com/connerdouglass/livestream-api/v1/utils"
	"github.com/gin-gonic/gin"
)

type StudioListStreamsReq struct {
	CreatorID uint64 `json:"creator_id"`
}

func StudioListStreams(
	creatorsService *services.CreatorsService,
	streamsService *services.StreamsService,
	membershipService *services.MembershipService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req StudioListStreamsReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the creator profile with the identifier
		creator, err := creatorsService.GetCreatorByID(req.CreatorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if creator == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "creator not found"})
			return
		}

		// Get the account from the context
		account := utils.CtxGetAccount(c)

		// Check if the account has access
		access, err := membershipService.IsMember(creator.ID, account.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !access {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		// Get all of the streams for the creator
		streams, err := streamsService.GetAllStreamsForCreatorID(creator.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Serialize all of the streams
		streamsSer := make([]map[string]interface{}, len(streams))
		for i := range streams {
			streamsSer[i] = serializeStreamForStudio(streams[i])
		}

		// Respond with the streams
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"streams": streamsSer,
			},
		})

	}
}
