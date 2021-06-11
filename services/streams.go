package services

import (
	"errors"
	"time"

	"github.com/godocompany/livestream-api/models"
	"github.com/godocompany/livestream-api/utils"
	"gorm.io/gorm"
)

// StreamsService manages the streams in the system
type StreamsService struct {
	DB *gorm.DB
}

type CreateStreamOptions struct {
	Title              string
	ScheduledStartDate time.Time
}

func (s *StreamsService) CreateStream(
	creator *models.CreatorProfile,
	options *CreateStreamOptions,
) (*models.Stream, error) {

	// Generate an identifier for the stream
	identifier, err := s.GenerateUnusedIdentifier()
	if err != nil {
		return nil, err
	}

	// Generate a stream key for the stream
	streamKey, err := s.GenerateUnusedStreamKey()
	if err != nil {
		return nil, err
	}

	// Create the stream
	stream := models.Stream{
		CreatorProfileID:   creator.ID,
		Identifier:         identifier,
		Title:              options.Title,
		StreamKey:          streamKey,
		Status:             models.StreamStatus_Upcoming,
		ScheduledStartDate: options.ScheduledStartDate,
		CreatedDate:        time.Now(),
	}
	if err := s.DB.Create(&stream).Error; err != nil {
		return nil, err
	}

	// Return the stream
	return &stream, nil

}

func (s *StreamsService) GetStreamByIdentifier(identifier string) (*models.Stream, error) {
	var stream models.Stream
	err := s.DB.
		Where("identifier = ?", identifier).
		Where("deleted_date IS NULL").
		First(&stream).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &stream, nil
}

func (s *StreamsService) GetStreamByStreamKey(streamKey string) (*models.Stream, error) {
	var stream models.Stream
	err := s.DB.
		Where("stream_key = ?", streamKey).
		Where("deleted_date IS NULL").
		First(&stream).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &stream, nil
}

func (s *StreamsService) GenerateUnusedIdentifier() (string, error) {
	maxAttempts := 1000
	for i := 0; i < maxAttempts; i++ {
		identifier := utils.RandHexStrInt64()
		stream, err := s.GetStreamByIdentifier(identifier)
		if err != nil {
			return "", err
		}
		if stream == nil {
			return identifier, nil
		}
	}
	return "", errors.New("GenerateUnusedIdentifier exceeded max attempts")
}

func (s *StreamsService) GenerateUnusedStreamKey() (string, error) {
	maxAttempts := 1000
	for i := 0; i < maxAttempts; i++ {
		streamKey := utils.RandHexStrInt64()
		stream, err := s.GetStreamByStreamKey(streamKey)
		if err != nil {
			return "", err
		}
		if stream == nil {
			return streamKey, nil
		}
	}
	return "", errors.New("GenerateUnusedStreamKey exceeded max attempts")
}
