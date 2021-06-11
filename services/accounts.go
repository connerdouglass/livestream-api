package services

import (
	"errors"

	"github.com/godocompany/livestream-api/models"
	"gorm.io/gorm"
)

// AccountsService manages the accounts on the platform
type AccountsService struct {
	DB *gorm.DB
}

// DoesAccountOwnStream checks if the given account owns the given stream
func (s *AccountsService) DoesAccountOwnStream(account *models.Account, stream *models.Stream) (bool, error) {

	// If either one is nil, return false
	if account == nil || stream == nil {
		return false, nil
	}

	// Check if we own a creator that owns the stream
	var creator models.CreatorProfile
	err := s.DB.
		Where("id = ?", stream.CreatorProfileID).
		Where("account_id = ?", account.ID).
		First(&creator).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, nil
	}
	return true, nil

}
