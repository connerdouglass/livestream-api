package hooks

import (
	"net/http"
	"time"

	"github.com/connerdouglass/livestream-api/services"
	"github.com/connerdouglass/livestream-api/v1/utils"
	"github.com/gin-gonic/gin"
)

type StudioCreateStreamReq struct {
	CreatorID uint64 `json:"creator_id"`
	Options   struct {
		Title              string `json:"title"`
		ScheduledStartDate int64  `json:"scheduled_start_date"`
	} `json:"options"`
}

func StudioCreateStream(
	creatorsService *services.CreatorsService,
	streamsService *services.StreamsService,
	membershipService *services.MembershipService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req StudioCreateStreamReq
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

		// Create a stream on the creator profile
		stream, err := streamsService.CreateStream(creator, &services.CreateStreamOptions{
			Title:              req.Options.Title,
			ScheduledStartDate: time.Unix(req.Options.ScheduledStartDate, 0),
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Return an empty response
		c.JSON(http.StatusOK, gin.H{
			"data": serializeStreamForStudio(stream),
		})

	}
}
