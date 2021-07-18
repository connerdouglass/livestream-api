package hooks

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/models"
	"github.com/godocompany/livestream-api/services"
	"github.com/godocompany/livestream-api/utils"
)

type GetCreatorMetaReq struct {
	Username string `json:"username"`
}

func GetCreatorMeta(
	creatorsService *services.CreatorsService,
	streamsService *services.StreamsService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req GetCreatorMetaReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get the creator with the given username
		creator, err := creatorsService.GetCreatorByUsername(req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if creator == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No such creator exists"})
			return
		}

		// Get the currently-live stream
		liveStream, err := streamsService.GetLiveStreamForCreator(creator)
		if err != nil {
			fmt.Println("Error fetching live stream: ", err.Error())
		}

		// Get the next upcoming stream (not yet live)
		nextStream, err := streamsService.GetNextStreamForCreator(creator)
		if err != nil {
			fmt.Println("Error fetching next upcoming stream: ", err.Error())
		}

		// Respond with the info about the creator
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"id":          creator.ID,
				"username":    creator.Username,
				"name":        creator.Name,
				"image":       creator.Image,
				"live_stream": serializeStream(liveStream),
				"next_stream": serializeStream(nextStream),
			},
		})

	}
}

func serializeStream(stream *models.Stream) map[string]interface{} {
	if stream == nil {
		return nil
	}
	return map[string]interface{}{
		"id":                   stream.ID,
		"identifier":           stream.Identifier,
		"title":                stream.Title,
		"status":               stream.Status,
		"scheduled_start_date": stream.ScheduledStartDate.Unix(),
		"current_viewers":      stream.CurrentViewers,
		"chatroom_url":         utils.FlattenNullString(stream.ChatRoomUrl),
	}
}
