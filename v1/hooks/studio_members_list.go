package hooks

import (
	"net/http"

	"github.com/connerdouglass/livestream-api/services"
	"github.com/connerdouglass/livestream-api/v1/utils"
	"github.com/gin-gonic/gin"
)

type StudioListMembersReq struct {
	CreatorID uint64 `json:"creator_id"`
}

func StudioListMembers(
	creatorsService *services.CreatorsService,
	membershipService *services.MembershipService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req StudioListMembersReq
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

		// Get all of the memberships
		members, err := membershipService.GetMembers(creator.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Serialize all of the members
		membersSer := make([]map[string]interface{}, len(members))
		for i, m := range members {
			membersSer[i] = map[string]interface{}{
				"id": m.ID,
				"account": map[string]interface{}{
					"id":    m.Account.ID,
					"email": m.Account.Email,
				},
			}
		}

		// Return an empty response
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{},
		})

	}
}
