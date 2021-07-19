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
	account, err := s.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	// Verify the password
	if !account.VerifyPassword(password) {
		return nil, nil
	}

	// Return the account
	return account, nil

}

// GetByEmail gets the account with the provided email address
func (s *AccountsService) GetByEmail(email string) (*models.Account, error) {
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
	return &account, nil
}
