package services

import (
	"database/sql"
	"errors"
	"fmt"
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

// GetAllStreamsForCreatorID gets all streams past, present, and future for the given creator ID
func (s *StreamsService) GetAllStreamsForCreatorID(creatorID uint64) ([]*models.Stream, error) {
	var streams []*models.Stream
	err := s.DB.
		Where("creator_profile_id = ?", creatorID).
		Where("deleted_date IS NULL").
		Find(&streams).
		Error
	if err != nil {
		return nil, err
	}
	return streams, nil
}

func (s *StreamsService) GetStreamByID(streamID uint64) (*models.Stream, error) {
	var stream models.Stream
	err := s.DB.
		Where("id = ?", streamID).
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

func (s *StreamsService) UpdateStreaming(stream *models.Stream, streaming bool) error {
	stream.Streaming = streaming
	if !streaming {
		stream.Status = models.StreamStatus_Ended
	}
	return s.DB.Save(stream).Error
}

func (s *StreamsService) UpdateStatus(stream *models.Stream, status string) error {

	// If the stream is nil
	if stream == nil {
		return errors.New("cannot update nil stream status")
	}

	// Define the slice of permitted status values
	possibleStatusValues := []string{
		models.StreamStatus_Upcoming,
		models.StreamStatus_Live,
		models.StreamStatus_Ended,
		models.StreamStatus_Cancelled,
	}

	// Make sure the status value is one of the permitted values
	contained := false
	for _, option := range possibleStatusValues {
		if option == status {
			contained = true
			break
		}
	}
	if !contained {
		return fmt.Errorf("unsupported stream status value: \"%s\"", status)
	}

	// Update the stream
	stream.Status = status
	if status == models.StreamStatus_Ended {
		stream.EndedDate = sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		}
	}
	return s.DB.Save(stream).Error

}

func (s *StreamsService) UpdateViewerCount(stream *models.Stream, count int) error {
	return s.DB.
		Model(&models.Stream{}).
		Where("id = ?", stream.ID).
		Update("current_viewers", count).
		Error
}

// GetLiveStreamForCreator gets the stream that is currently live for a creator
func (s *StreamsService) GetLiveStreamForCreator(creator *models.CreatorProfile) (*models.Stream, error) {
	var stream models.Stream
	err := s.DB.
		Where("creator_profile_id = ?", creator.ID).
		Where("status = ?", models.StreamStatus_Live).
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

// GetNextStreamForCreator gets the stream that is next upcoming (and not currently live) for a creator
func (s *StreamsService) GetNextStreamForCreator(creator *models.CreatorProfile) (*models.Stream, error) {
	var stream models.Stream
	err := s.DB.
		Where("creator_profile_id = ?", creator.ID).
		Where("status = ?", models.StreamStatus_Upcoming).
		Order("scheduled_start_date ASC").
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
