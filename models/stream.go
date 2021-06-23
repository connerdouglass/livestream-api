package models

import (
	"database/sql"
	"time"
)

const (
	StreamStatus_Upcoming  = "upcoming"
	StreamStatus_Live      = "live"
	StreamStatus_Ended     = "ended"
	StreamStatus_Cancelled = "cancelled"
)

// Stream represents a scheduled or currently-live stream
type Stream struct {
	ID                 uint64 `gorm:"primaryKey"`
	CreatorProfileID   uint64
	CreatorProfile     *CreatorProfile
	Identifier         string
	Title              string
	StreamKey          string
	Status             string
	Streaming          bool
	ScheduledStartDate time.Time
	CurrentViewers     int
	EndedDate          sql.NullTime
	CreatedDate        time.Time
	DeletedDate        sql.NullTime
}
