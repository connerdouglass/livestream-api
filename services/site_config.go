package services

import (
	"errors"

	"github.com/connerdouglass/livestream-api/models"
	"gorm.io/gorm"
)

// SiteConfigService manages the accounts on the platform
type SiteConfigService struct {
	DB *gorm.DB
}

// GetSiteConfig gets the site config from the database, or creates a new one if needed
func (s *SiteConfigService) GetSiteConfig() (*models.SiteConfig, error) {

	// Get the site configuration
	var config models.SiteConfig
	err := s.DB.First(&config).Error
	if err == nil {
		return &config, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create a new configuration
	config = models.SiteConfig{}
	if err := s.DB.Create(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil

}
