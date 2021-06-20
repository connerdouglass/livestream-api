package hooks

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/services"
)

type AuthLoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AuthLogin(
	accountsService *services.AccountsService,
	authTokensService *services.AuthTokensService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the request body
		var req AuthLoginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Find the account with the provided email and password
		account, err := accountsService.FindByLogin(
			req.Email,
			req.Password,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if account == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect email or password"})
			return
		}

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
