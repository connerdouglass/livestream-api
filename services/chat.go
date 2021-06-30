package services

import (
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
