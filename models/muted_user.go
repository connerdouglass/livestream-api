package models

import (
	"database/sql"
	"time"
)

// MutedUser is a user that is muted in chat
type MutedUser struct {
	ID               uint64 `gorm:"primaryKey"`
	CreatorProfileID uint64
	CreatorProfile   *CreatorProfile
	Username         string
	CreatedDate      time.Time
	DeletedDate      sql.NullTime
}
