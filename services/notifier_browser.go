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

type BrowserNotifier struct {
	DB                *gorm.DB
	SiteConfigService *SiteConfigService
}

func (bn *BrowserNotifier) NotifySubscribers(
	creatorID uint64,
	notification *Notification,
) error {

	// Get all of the browser notify targets
	var targets []*models.BrowserNotifyTarget
	err := bn.DB.
		Where("deleted_date IS NULL").
		Where(
			"id IN (?)",
			bn.DB.
				Select("browser_notify_target_id").
				Model(&models.TelegramNotifySub{}).
				Where("deleted_date IS NULL").
				Where("creator_profile_id = ?", creatorID),
		).
		Find(&targets).
		Error
	if err != nil {
		return nil
	}

	// Format the title and body as JSON
	message, err := json.Marshal(map[string]interface{}{
		"title": notification.Title,
		"body":  notification.Body,
		"link":  notification.Link,
		"image": notification.Image,
	})
	if err != nil {
		return err
	}

	// Create the wait group to run all tasks in parallel
	var wg sync.WaitGroup
	wg.Add(len(targets))

	// Loop through all of the targets
	for index := range targets {
		go func(i int) {

			// Defer a cleanup function
			defer wg.Done()

			// Send the notification
			err := bn.SendBrowserNotification(
				targets[i].RegistrationData,
				message,
			)
			if err != nil {
				fmt.Println("Error sending browser notification: ", err.Error())
			}

		}(index)
	}

	// Wait for all tasks to complete
	wg.Wait()

	// Return without error
	return nil

}

func (bn *BrowserNotifier) SendBrowserNotification(
	registrationData string,
	message []byte,
) error {

	// Decode subscription
	sub := webpush.Subscription{}
	if err := json.Unmarshal([]byte(registrationData), &sub); err != nil {
		return err
	}

	// Get the vapid keypair
	keypair, err := bn.GetVapidKeyPair()
	if err != nil {
		return err
	}

	// Send the browser push notification
	resp, err := webpush.SendNotification(message, &sub, &webpush.Options{
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

// GetVapidKeyPair gets the keypair for VAPID keys
func (bn *BrowserNotifier) GetVapidKeyPair() (*VapidKeyPair, error) {

	// Get the site config
	config, err := bn.SiteConfigService.GetSiteConfig()
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
	if err := bn.DB.Save(config).Error; err != nil {
		return nil, err
	}

	// Return the keypair object
	return &VapidKeyPair{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}, nil

}

func (bn *BrowserNotifier) getOrCreateTarget(regData string) (*models.BrowserNotifyTarget, error) {

	// Get the notify target with this registration data
	var target models.BrowserNotifyTarget
	err := bn.DB.
		Where("deleted_date IS NULL").
		Where("registration_data = ?", regData).
		First(&target).
		Error

	// If the target was found, return it
	if err == nil {
		return &target, nil
	}

	// If the error is something other than "not found"
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create the new target
	target = models.BrowserNotifyTarget{
		RegistrationData: regData,
		CreatedDate:      time.Now(),
	}
	if err := bn.DB.Save(&target).Error; err != nil {
		return nil, err
	}

	// Return the new target instance
	return &target, nil

}

func (bn *BrowserNotifier) getNotifySub(
	targetID uint64,
	creatorID uint64,
) (*models.BrowserNotifySub, error) {
	var sub models.BrowserNotifySub
	err := bn.DB.
		Where("deleted_date IS NULL").
		Where("browser_notify_target_id = ?", targetID).
		Where("creator_profile_id = ?", creatorID).
		First(&sub).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

func (bn *BrowserNotifier) UpdateSub(
	registrationData string,
	creatorID uint64,
	subscribed bool,
) error {

	// Get or create the target
	target, err := bn.getOrCreateTarget(registrationData)
	if err != nil {
		return err
	}

	// Get the notification subscription
	sub, err := bn.getNotifySub(target.ID, creatorID)
	if err != nil {
		return err
	}

	// If our job is already done, return here
	if (sub == nil) == !subscribed {
		return nil
	}

	// If we're un-subscribing
	if !subscribed {
		sub.DeletedDate = sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		}
		return bn.DB.Save(sub).Error
	}

	// Create a new subscription
	sub = &models.BrowserNotifySub{
		BrowserNotifyTargetID: target.ID,
		CreatorProfileID:      creatorID,
		CreatedDate:           time.Now(),
	}
	return bn.DB.Create(sub).Error

}

func (bn *BrowserNotifier) RegisterTarget(registrationData string) error {
	_, err := bn.getOrCreateTarget(registrationData)
	return err
}

func (bn *BrowserNotifier) GetAllSubs(registrationData string) ([]*models.BrowserNotifySub, error) {

	// Get or create the target
	target, err := bn.getOrCreateTarget(registrationData)
	if err != nil {
		return nil, err
	}

	// Get all of the browser subscriptions
	var subs []*models.BrowserNotifySub
	err = bn.DB.
		Where("deleted_date IS NULL").
		Where("browser_notify_target_id = ?", target.ID).
		Find(&subs).
		Error
	if err != nil {
		return nil, err
	}

	// Return the slice of subscribers
	return subs, nil

}
