package services

import (
	"errors"
	"time"

	"github.com/godocompany/livestream-api/models"
	"gorm.io/gorm"
)

// ChatService manages chat moderation
type ChatService struct {
	DB *gorm.DB
}

func (s *ChatService) MuteUser( /*creatorID uint64,*/ username string) (*models.MutedUser, error) {

	// Add an entry to mute the user
	mutedUser := models.MutedUser{
		Username:    username,
		CreatedDate: time.Now(),
	}
	if err := s.DB.Create(&mutedUser).Error; err != nil {
		return nil, err
	}
	return &mutedUser, nil

}

func (s *ChatService) IsUserMuted( /* creatorID uint64,*/ username string) (bool, error) {
	var mutedUser models.MutedUser
	err := s.DB.
		Where("deleted_date IS NULL").
		Where("username LIKE ?", username).
		// Where("creator_profile_id = ?", creatorID).
		First(&mutedUser).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
