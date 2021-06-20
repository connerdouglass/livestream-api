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

// FindByLogin finds an account with the provided login credentials
func (s *AccountsService) FindByLogin(email, password string) (*models.Account, error) {

	// Find the account with the email
	var account models.Account
	err := s.DB.
		Where("email LIKE ?", email).
		First(&account).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Verify the password
	if !account.VerifyPassword(password) {
		return nil, nil
	}

	// Return the account
	return &account, nil

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
