package services

import (
	"errors"
	"regexp"
	"strings"

	"github.com/godocompany/livestream-api/models"
	"gorm.io/gorm"
)

// CreatorsService manages the creators on the platform
type CreatorsService struct {
	DB *gorm.DB
}

// GetCreatorByUsername gets the creator with the provided username
func (s *CreatorsService) GetCreatorByUsername(username string) (*models.CreatorProfile, error) {

	// Trim the username
	username = strings.TrimSpace(username)

	// If the username is invalid
	if !s.ValidateUsername(username) {
		return nil, nil
	}

	// Query for theu creator
	var creator models.CreatorProfile
	err := s.DB.
		Where("deleted_date IS NULL").
		Where("username LIKE ?", username).
		First(&creator).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &creator, nil

}

// DoesCreatorOwnStream checks if the given creator owns the given stream
func (s *CreatorsService) DoesCreatorOwnStream(creator *models.CreatorProfile, stream *models.Stream) bool {

	// If either one is nil, return false
	if creator == nil || stream == nil {
		return false
	}

	// Check if the identifiers match
	return creator.ID == stream.CreatorProfileID

}

// GetCreatorByID gets the creator with the given identifier
func (s *CreatorsService) GetCreatorByID(creatorID uint64) (*models.CreatorProfile, error) {
	var creator models.CreatorProfile
	err := s.DB.
		Where("deleted_date IS NULL").
		Where("id = ?", creatorID).
		First(&creator).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &creator, nil
}

// ValidateUsernae checks if the provided username is valid
func (s *CreatorsService) ValidateUsername(username string) bool {
	pattern := regexp.MustCompile(`^\w+$`)
	return pattern.MatchString(username)
}

// GetCreatorsByAccountID gets all the creator profiles belonging to an account
func (s *CreatorsService) GetCreatorsByAccountID(accountID uint64) ([]*models.CreatorProfile, error) {
	var creators []*models.CreatorProfile
	err := s.DB.
		Where("account_id = ?", accountID).
		Where("deleted_date IS NULL").
		Find(creators).
		Error
	if err != nil {
		return nil, err
	}
	return creators, nil
}
