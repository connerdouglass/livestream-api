package models

import (
	"database/sql"
	"time"
)

// NotificationSubscriber represents a subscriber to browser notifications
type NotificationSubscriber struct {
	ID               uint64 `gorm:"primaryKey"`
	CreatorProfileID sql.NullInt64
	CreatorProfile   *CreatorProfile
	RegistrationData sql.NullString
	TelegramChatID   sql.NullInt64
	CreatedDate      time.Time
	DeletedDate      sql.NullTime
}
