package hooks

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/services"
	"github.com/godocompany/livestream-api/v1/utils"
)

type StudioCreateStreamReq struct {
	CreatorID uint64 `json:"creator_id"`
	Options   struct {
		Title              string `json:"title"`
		ScheduledStartDate int64  `json:"scheduled_start_date"`
	} `json:"options"`
}

func StudioCreateStream(
	accountsService *services.AccountsService,
	creatorsService *services.CreatorsService,
	streamsService *services.StreamsService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req StudioCreateStreamReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the account from the context
		account := utils.CtxGetAccount(c)

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
		if creator.AccountID != account.ID {
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
