package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/godocompany/livestream-api/models"
	"gorm.io/gorm"
)

type VapidKeyPair struct {
	PublicKey  string
	PrivateKey string
}

// NotificationsService manages push notifications
type NotificationsService struct {
	DB                *gorm.DB
	SiteConfigService *SiteConfigService
	TelegramService   *TelegramService
}

// GetVapidKeyPair gets the keypair for VAPID keys
func (s *NotificationsService) GetVapidKeyPair() (*VapidKeyPair, error) {

	// Get the site config
	config, err := s.SiteConfigService.GetSiteConfig()
	if err != nil {
		return nil, err
	}

	// If there are already keys in the config
	if config.VapidPublicKey.Valid && config.VapidPrivateKey.Valid {
		return &VapidKeyPair{
			PublicKey:  config.VapidPublicKey.String,
			PrivateKey: config.VapidPrivateKey.String,
		}, nil
	}

	// Generate the keys
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		return nil, err
	}

	// Store the keys in the database
	config.VapidPublicKey = sql.NullString{
		Valid:  true,
		String: publicKey,
	}
	config.VapidPrivateKey = sql.NullString{
		Valid:  true,
		String: privateKey,
	}
	if err := s.DB.Save(config).Error; err != nil {
		return nil, err
	}

	// Return the keypair object
	return &VapidKeyPair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil

}

func (s *NotificationsService) Subscribe(
	creatorID uint64,
	browserRegistrationData *string,
	telegramChatID *int64,
) error {

	// Construct the query
	query := s.DB.
		Model(&models.NotificationSubscriber{}).
		Where("deleted_date IS NULL")

	// If it's a browser registration
	if browserRegistrationData != nil {
		query = query.Where("registration_data = ?", *browserRegistrationData)
	} else if telegramChatID != nil {
		query = query.Where("telegram_chat_id = ?", *telegramChatID)
	} else {
		return errors.New("cannot subscribe without browser or telegram source")
	}

	// Check if there is already a subscription
	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	// Create the subscription
	sub := models.NotificationSubscriber{
		CreatorProfileID: sql.NullInt64{
			Valid: true,
			Int64: int64(creatorID),
		},
		CreatedDate: time.Now(),
	}
	if browserRegistrationData != nil {
		sub.RegistrationData = sql.NullString{
			Valid:  true,
			String: *browserRegistrationData,
		}
	} else if telegramChatID != nil {
		sub.TelegramChatID = sql.NullInt64{
			Valid: true,
			Int64: *telegramChatID,
		}
	}
	return s.DB.Create(&sub).Error

}

type BrowserNotification struct {
	RegistrationData string
	Message          []byte
}

func (s *NotificationsService) SendBrowserNotification(options *BrowserNotification) error {

	// Decode subscription
	sub := webpush.Subscription{}
	if err := json.Unmarshal([]byte(options.RegistrationData), &sub); err != nil {
		return err
	}

	// Get the vapid keypair
	keypair, err := s.GetVapidKeyPair()
	if err != nil {
		return err
	}

	// Send the browser push notification
	resp, err := webpush.SendNotification(options.Message, &sub, &webpush.Options{
		// Subscriber:      "example@example.com",
		VAPIDPublicKey:  keypair.PublicKey,
		VAPIDPrivateKey: keypair.PrivateKey,
		TTL:             30,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Return without error
	return nil

}

type TelegramNotification struct {
	ChatID  int64
	Message string
}

func (s *NotificationsService) SendTelegramNotification(options *TelegramNotification) error {

	// Send the message over Telegram
	return s.TelegramService.SendMessage(options.ChatID, options.Message)

}

func (s *NotificationsService) SendNotificationToCreatorSubscribers(
	creatorID uint64,
	title string,
	body string,
) error {

	// Format the title and body as JSON
	message, err := json.Marshal(map[string]string{
		"title": title,
		"body":  body,
	})
	if err != nil {
		return err
	}

	// Get all of the recipients
	var subscribers []*models.NotificationSubscriber
	err = s.DB.
		Where("deleted_date IS NULL").
		Where("creator_profile_id = ?", creatorID).
		Find(&subscribers).
		Error
	if err != nil {
		return err
	}

	// Seperate out the browser notifications and telegram notifications
	var browserSubscribers []string
	var telegramChatIDs []int64
	for _, sub := range subscribers {
		if sub.RegistrationData.Valid {
			browserSubscribers = append(browserSubscribers, sub.RegistrationData.String)
		} else if sub.TelegramChatID.Valid {
			telegramChatIDs = append(telegramChatIDs, sub.TelegramChatID.Int64)
		}
	}

	// Create a wait group for all the subscribers
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		s.bulkTelegramNotify(telegramChatIDs, fmt.Sprintf("%s: %s", title, body))
	}()

	go func() {
		defer wg.Done()
		s.bulkBrowserNotify(browserSubscribers, message)
	}()

	// Wait for all the notifications to finish
	wg.Wait()

	// Return without error
	return nil

}

func (s *NotificationsService) bulkBrowserNotify(
	registrationDatas []string,
	message []byte,
) {

	// Create a wait group for all the subscribers
	var wg sync.WaitGroup
	wg.Add(len(registrationDatas))

	// Loop through the notifications
	for _, regData := range registrationDatas {
		go func(registrationData string) {

			// Defer the cleanup
			defer wg.Done()

			// Send the browser notification
			err := s.SendBrowserNotification(&BrowserNotification{
				RegistrationData: registrationData,
				Message:          message,
			})
			if err != nil {
				fmt.Println("Error sending notification: ", err.Error())
			}

		}(regData)
	}

	// Wait for them all to complete
	wg.Wait()

}

func (s *NotificationsService) bulkTelegramNotify(
	chatIDs []int64,
	message string,
) {

	// Create a wait group for all the subscribers
	var wg sync.WaitGroup
	wg.Add(len(chatIDs))

	// Loop through the telegram chat ids
	for _, telegramChatID := range chatIDs {
		go func(chatID int64) {

			// Defer the cleanup
			defer wg.Done()

			// Send the notification
			err := s.SendTelegramNotification(&TelegramNotification{
				ChatID:  chatID,
				Message: message,
			})
			if err != nil {
				fmt.Println("Error sending notification: ", err.Error())
			}

		}(telegramChatID)
	}

	// Wait for all the notifications to finish
	wg.Wait()

}
