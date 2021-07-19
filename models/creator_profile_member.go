package models

import (
	"database/sql"
	"time"
)

// CreatorProfileMember is a member with access to manage a profile
type CreatorProfileMember struct {
	ID               uint64 `gorm:"primaryKey"`
	AccountID        uint64
	Account          *Account
	CreatorProfileID uint64
	CreatorProfile   *CreatorProfile
	CreatedDate      time.Time
	DeletedDate      sql.NullTime
}
