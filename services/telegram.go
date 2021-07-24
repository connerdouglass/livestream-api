package services

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/godocompany/livestream-api/models"
	"github.com/godocompany/livestream-api/utils"
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

func (s *TelegramService) Subscribe(
	subscriberData *TelegramNotificationSubscriber,
) error {

	// Check if there is already a subscription
	var count int64
	err := s.DB.
		Model(&models.NotificationSubscriber{}).
		Where("deleted_date IS NULL").
		Where("telegram_chat_id = ?", subscriberData.ChatID).
		Count(&count).
		Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// Create the subscription
	sub := models.NotificationSubscriber{
		TelegramUserID: sql.NullInt64{
			Valid: true,
			Int64: subscriberData.UserID,
		},
		TelegramChatID: sql.NullInt64{
			Valid: true,
			Int64: subscriberData.ChatID,
		},
		CreatedDate: time.Now(),
	}
	return s.DB.Create(&sub).Error

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

			// Create the subscriber data
			subData := TelegramNotificationSubscriber{
				UserID: int64(update.Message.From.ID),
				ChatID: update.Message.Chat.ID,
			}

			// Subscribe to notifications
			if err := s.Subscribe(&subData); err != nil {
				fmt.Println("Error subscribing to Telegram notifications: ", err.Error())
			}

		}

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID

		// bot.Send(msg)

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
