package models

import (
	"database/sql"
	"time"
)

// BrowserNotifySub represents a subscription to receive browser notifications from a specific creator profile
type BrowserNotifySub struct {
	ID                    uint64 `gorm:"primaryKey"`
	BrowserNotifyTargetID uint64
	BrowserNotifyTarget   *BrowserNotifyTarget
	CreatorProfileID      uint64
	CreatorProfile        *CreatorProfile
	CreatedDate           time.Time
	DeletedDate           sql.NullTime
}
