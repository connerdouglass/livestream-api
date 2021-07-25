package models

import (
	"database/sql"
	"time"
)

// TelegramNotifyTarget represents a Telegram user who receives notifications
type TelegramNotifyTarget struct {
	ID             uint64 `gorm:"primaryKey"`
	TelegramChatID int64
	TelegramUserID int64
	CreatedDate    time.Time
	DeletedDate    sql.NullTime
}
