package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/connerdouglass/livestream-api/models"
	"github.com/connerdouglass/livestream-api/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gorm.io/gorm"
)

type TelegramUser struct {
	AuthDate  uint64  `json:"auth_date"`
	ID        uint64  `json:"id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	PhotoUrl  *string `json:"photo_url"`
	Username  string  `json:"username"`
	Hash      string  `json:"hash"`
}

type TelegramNotificationSubscriber struct {
	ChatID int64
	UserID int64
}

type TelegramService struct {
	DB          *gorm.DB
	BotAPIKey   string
	BotUsername string
}

// Verify verified the validity of a Telegram user
func (s *TelegramService) Verify(user *TelegramUser) bool {

	// Add all of the keyval pairs to a slice
	keyvals := []string{}
	keyvals = append(keyvals, fmt.Sprintf("auth_date=%d", user.AuthDate))
	if user.FirstName != nil {
		keyvals = append(keyvals, fmt.Sprintf("first_name=%s", *user.FirstName))
	}
	// keyvals = append(keyvals, fmt.Sprintf("hash=%s", user.Hash))
	keyvals = append(keyvals, fmt.Sprintf("id=%d", user.ID))
	if user.LastName != nil {
		keyvals = append(keyvals, fmt.Sprintf("last_name=%s", *user.LastName))
	}
	if user.PhotoUrl != nil {
		keyvals = append(keyvals, fmt.Sprintf("photo_url=%s", *user.PhotoUrl))
	}
	keyvals = append(keyvals, fmt.Sprintf("username=%s", user.Username))

	// Sort the keyvals slice
	// sort.Strings(keyvals)

	// Create the data check string
	dataCheckString := strings.Join(keyvals, "\n")

	// Create the hash of the API key
	secretKey := utils.Sha256Hex(s.BotAPIKey)

	// Calculate the actual hash
	hash := utils.HmacSha256(dataCheckString, secretKey)

	// Compare the hashes
	return hash == user.Hash

}

func (s *TelegramService) RegisterTarget(
	userID int64,
	chatID int64,
) error {

	// Check if there is already a notification target for this user
	var count int64
	err := s.DB.
		Model(&models.TelegramNotifyTarget{}).
		Where("deleted_date IS NULL").
		Where("telegram_chat_id = ?", chatID).
		Where("telegram_user_id = ?", userID).
		Count(&count).
		Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// Create the notification target
	target := models.TelegramNotifyTarget{
		TelegramUserID: userID,
		TelegramChatID: chatID,
		CreatedDate:    time.Now(),
	}
	return s.DB.Create(&target).Error

}

func (s *TelegramService) Listen() error {

	// Create the Telegram bot client
	bot, err := tgbotapi.NewBotAPI(s.BotAPIKey)
	if err != nil {
		return err
	}
	bot.Debug = true

	// Create a channel for bot updates
	updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Offset:  0,
		Limit:   0,
		Timeout: 60,
	})
	if err != nil {
		return err
	}

	// Iterate through the updates on the channel
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// If the message is /start
		if update.Message.Text == "/start" {

			// Subscribe to notifications
			if err := s.RegisterTarget(
				int64(update.Message.From.ID),
				update.Message.Chat.ID,
			); err != nil {
				fmt.Println("Error subscribing to Telegram notifications: ", err.Error())
			}

		}

	}

	// Return no error
	return nil

}

func (s *TelegramService) SendMessage(chatID int64, message string) error {

	// Create the Telegram bot client
	bot, err := tgbotapi.NewBotAPI(s.BotAPIKey)
	if err != nil {
		return err
	}
	bot.Debug = true

	// Create and send the message
	msg := tgbotapi.NewMessage(chatID, message)

	// Send the message
	if _, err := bot.Send(msg); err != nil {
		return err
	}

	// Return without error
	return nil
}
