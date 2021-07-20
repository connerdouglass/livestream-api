package models

import (
	"database/sql"
	"time"
)

// NotificationSubscriber represents a subscriber to browser notifications
type NotificationSubscriber struct {
	ID               uint64 `gorm:"primaryKey"`
	CreatorProfileID uint64
	CreatorProfile   *CreatorProfile
	RegistrationData sql.NullString
	CreatedDate      time.Time
	DeletedDate      sql.NullTime
}
