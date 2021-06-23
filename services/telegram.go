package services

import (
	"fmt"
	"strings"

	"github.com/godocompany/livestream-api/utils"
)

type TelegramUser struct {
	AuthDate  uint64 `json:"auth_date"`
	ID        uint64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Hash      string `json:"hash"`
}

type TelegramService struct {
	BotAPIKey   string
	BotUsername string
}

// Verify verified the validity of a Telegram user
func (s *TelegramService) Verify(user *TelegramUser) bool {

	// Create the data check string
	dataCheckString := strings.Join(
		[]string{
			fmt.Sprintf("auth_date=%d", user.AuthDate),
			fmt.Sprintf("first_name=%s", user.FirstName),
			fmt.Sprintf("hash=%s", user.Hash),
			fmt.Sprintf("id=%d", user.ID),
			fmt.Sprintf("last_name=%s", user.LastName),
			fmt.Sprintf("username=%s", user.Username),
		},
		"\n",
	)

	// Create the hash of the API key
	secretKey := utils.Sha256Hex(s.BotAPIKey)

	// Calculate the actual hash
	hash := utils.HmacSha256(dataCheckString, secretKey)

	// Compare the hashes
	return hash == user.Hash

}
