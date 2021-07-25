package models

import (
	"database/sql"
	"time"
)

// BrowserNotifyTarget represents a user who can receive browser notifications
type BrowserNotifyTarget struct {
	ID               uint64 `gorm:"primaryKey"`
	RegistrationData string
	CreatedDate      time.Time
	DeletedDate      sql.NullTime
}
