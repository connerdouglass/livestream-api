package hooks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/services"
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

		// Respond with the info about the creator
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"id":       creator.ID,
				"username": creator.Username,
				"name":     creator.Name,
			},
		})

	}
}
