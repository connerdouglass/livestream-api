package models

import (
	"database/sql"
	"time"
)

// CreatorProfile is a profile on the platform
type CreatorProfile struct {
	ID          uint64 `gorm:"primaryKey"`
	AccountID   uint64
	Account     *Account
	Username    string
	Name        string
	CreatedDate time.Time
	DeletedDate sql.NullTime
}
