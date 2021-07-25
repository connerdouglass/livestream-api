package services

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/godocompany/livestream-api/models"
	"gorm.io/gorm"
)

type TelegramNotifier struct {
	DB              *gorm.DB
	TelegramService *TelegramService
}

func (tn *TelegramNotifier) NotifySubscribers(
	creatorID uint64,
	notification *Notification,
) error {

	// Get all of the telegram notify targets
	var targets []*models.TelegramNotifyTarget
	err := tn.DB.
		Where("deleted_date IS NULL").
		Where(
			"id IN (?)",
			tn.DB.
				Select("telegram_notify_target_id").
				Model(&models.TelegramNotifySub{}).
				Where("deleted_date IS NULL").
				Where("creator_profile_id = ?", creatorID),
		).
		Find(&targets).
		Error
	if err != nil {
		return nil
	}

	// Format the message parts into one string to send via Telegram
	message := fmt.Sprintf("%s: %s", notification.Title, notification.Body)
	if notification.Link != nil {
		message = fmt.Sprintf("%s\n\n%s", message, *notification.Link)
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
			err := tn.TelegramService.SendMessage(
				targets[i].TelegramChatID,
				message,
			)
			if err != nil {
				fmt.Println("Error sending Telegram message: ", err.Error())
			}

		}(index)
	}

	// Wait for all tasks to complete
	wg.Wait()

	// Return without error
	return nil

}

func (tn *TelegramNotifier) getNotifyTarget(user *TelegramUser) (*models.TelegramNotifyTarget, error) {
	var target models.TelegramNotifyTarget
	err := tn.DB.
		Where("deleted_date IS NULL").
		Where("telegram_user_id = ?", user.ID).
		First(&target).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &target, nil
}

func (tn *TelegramNotifier) getNotifySub(
	targetID uint64,
	creatorID uint64,
) (*models.TelegramNotifySub, error) {
	var sub models.TelegramNotifySub
	err := tn.DB.
		Where("deleted_date IS NULL").
		Where("telegram_notify_target_id = ?", targetID).
		Where("creator_profile_id = ?", creatorID).
		First(&sub).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, nil
	}
	return &sub, nil
}

func (tn *TelegramNotifier) UpdateSub(
	user *TelegramUser,
	creatorID uint64,
	subscribed bool,
) error {

	// Get or create the target
	target, err := tn.getNotifyTarget(user)
	if err != nil {
		return err
	}
	if target == nil {
		return errors.New("user is not registered for notifications")
	}

	// Get the notification subscription
	sub, err := tn.getNotifySub(target.ID, creatorID)
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
		return tn.DB.Save(sub).Error
	}

	// Create a new subscription
	sub = &models.TelegramNotifySub{
		TelegramNotifyTargetID: target.ID,
		CreatorProfileID:       creatorID,
		CreatedDate:            time.Now(),
	}
	return tn.DB.Create(sub).Error

}

func (tn *TelegramNotifier) GetAllSubs(user *TelegramUser) (bool, []*models.TelegramNotifySub, error) {

	// Get or create the target
	target, err := tn.getNotifyTarget(user)
	if err != nil {
		return false, nil, err
	}

	// Get all of the telegram subscriptions
	var subs []*models.TelegramNotifySub
	err = tn.DB.
		Where("deleted_date IS NULL").
		Where("telegram_notify_target_id = ?", target.ID).
		Find(&subs).
		Error
	if err != nil {
		return true, nil, err
	}

	// Return the slice of subscribers
	return true, subs, nil

}
