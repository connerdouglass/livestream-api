package hooks

import (
	"errors"
	"net/http"
	"time"

	"github.com/connerdouglass/livestream-api/models"
	"github.com/connerdouglass/livestream-api/services"
	"github.com/connerdouglass/livestream-api/v1/utils"
	"github.com/gin-gonic/gin"
)

func AuthWhoAmI(
	authTokensService *services.AuthTokensService,
	membershipService *services.MembershipService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the account from the request
		account := utils.CtxGetAccount(c)

		// Serialize the whoami info
		whoami, err := serializeWhoAmI(
			account,
			authTokensService,
			membershipService,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the whoami info for this account
		c.JSON(http.StatusOK, gin.H{
			"data": whoami,
		})

	}
}

func serializeWhoAmI(
	account *models.Account,
	authTokensService *services.AuthTokensService,
	membershipService *services.MembershipService,
) (map[string]interface{}, error) {

	// Return nil if the account is nil
	if account == nil {
		return nil, errors.New("something went wrong")
	}

	// Create an authentication token for the account
	token, err := authTokensService.CreateToken(
		account,
		time.Now(),
		time.Now().Add(time.Hour*24*30),
	)
	if err != nil {
		return nil, err
	}

	// Get all of the creator profiles on this account
	creators, err := membershipService.GetCreatorProfiles(account.ID)
	if err != nil {
		return nil, err
	}

	// Serialize all of the creators
	creatorsSer := make([]map[string]interface{}, len(creators))
	for i, creator := range creators {
		creatorsSer[i] = map[string]interface{}{
			"id":       creator.ID,
			"username": creator.Username,
			"name":     creator.Name,
		}
	}

	// Return the map of whoami info
	return map[string]interface{}{
		"id":       account.ID,
		"email":    account.Email,
		"token":    token,
		"creators": creatorsSer,
	}, nil
}
