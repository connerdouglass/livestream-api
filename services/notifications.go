package services

import (
	"database/sql"
	"encoding/json"
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
		PublicKey:  config.VapidPublicKey.String,
		PrivateKey: config.VapidPrivateKey.String,
	}, nil

}

func (s *NotificationsService) Subscribe(
	creatorID uint64,
	registrationData string,
) error {
	sub := models.NotificationSubscriber{
		CreatorProfileID: creatorID,
		RegistrationData: sql.NullString{
			Valid:  true,
			String: registrationData,
		},
		CreatedDate: time.Now(),
	}
	return s.DB.Create(&sub).Error
}

type BrowserNotification struct {
	RegistrationData string
	Message          []byte
}

func (s *NotificationsService) SendBrowserNotification(options *BrowserNotification) error {

	// Decode subscription
	sub := &webpush.Subscription{}
	if err := json.Unmarshal([]byte(options.RegistrationData), s); err != nil {
		return err
	}

	// Get the vapid keypair
	keypair, err := s.GetVapidKeyPair()
	if err != nil {
		return err
	}

	// Send the browser push notification
	resp, err := webpush.SendNotification(options.Message, sub, &webpush.Options{
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

func (s *NotificationsService) SendNotificationToCreatorSubscribers(
	creatorID uint64,
	message []byte,
) error {

	// Get all of the recipients
	var subscribers []*models.NotificationSubscriber
	err := s.DB.
		Where("deleted_date IS NULL").
		Where("creator_profile_id = ?", creatorID).
		Find(&subscribers).
		Error
	if err != nil {
		return err
	}

	// Seperate out the browser notifications and telegram notifications
	var browserSubscribers []string
	for _, sub := range subscribers {
		if sub.RegistrationData.Valid {
			browserSubscribers = append(browserSubscribers, sub.RegistrationData.String)
		} else {
			// TODO: Handle Telegram notifications
		}
	}

	// Create a wait group for all the subscribers
	var wg sync.WaitGroup
	wg.Add(len(browserSubscribers))

	// Loop through the notifications
	for _, regData := range browserSubscribers {
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

	// Wait for all the notifications to finish
	wg.Done()

	// Return without error
	return nil

}
