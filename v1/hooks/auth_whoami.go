package hooks

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/services"
	"github.com/godocompany/livestream-api/v1/utils"
)

func AuthWhoAmI(
	authTokensService *services.AuthTokensService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the account from the request
		account := utils.CtxGetAccount(c)

		// Create an authentication token for the account
		token, err := authTokensService.CreateToken(
			account,
			time.Now(),
			time.Now().Add(time.Hour*24*30),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the whoami info for this account
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"id":    account.ID,
				"email": account.Email,
				"token": token,
			},
		})

	}
}
