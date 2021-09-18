package hooks

import (
	"net/http"

	"github.com/connerdouglass/livestream-api/services"
	"github.com/connerdouglass/livestream-api/v1/utils"
	"github.com/gin-gonic/gin"
)

type StudioAddMemberReq struct {
	CreatorID uint64 `json:"creator_id"`
	Email     string `json:"email"`
}

func StudioAddMember(
	accountsService *services.AccountsService,
	creatorsService *services.CreatorsService,
	membershipService *services.MembershipService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req StudioAddMemberReq
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

		// Find the target account with the email address
		targetAccount, err := accountsService.GetByEmail(req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if targetAccount == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "target account not found"})
			return
		}

		// Create the membership
		if err := membershipService.AddMember(creator.ID, targetAccount.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return an empty response
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{},
		})

	}
}
