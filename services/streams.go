package services

import "gorm.io/gorm"

// StreamsService manages the streams in the system
type StreamsService struct {
	DB *gorm.DB
}
