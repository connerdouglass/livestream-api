package models

import (
	"database/sql"
	"time"
)

// TelegramNotifySub represents a subscription to receive Telegram notifications for a specific creator profile
type TelegramNotifySub struct {
	ID                     uint64 `gorm:"primaryKey"`
	TelegramNotifyTargetID uint64
	TelegramNotifyTarget   *TelegramNotifyTarget
	CreatorProfileID       uint64
	CreatorProfile         *CreatorProfile
	CreatedDate            time.Time
	DeletedDate            sql.NullTime
}
