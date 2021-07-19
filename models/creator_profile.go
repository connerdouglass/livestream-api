package models

import (
	"database/sql"
	"time"
)

// CreatorProfile is a profile on the platform
type CreatorProfile struct {
	ID          uint64 `gorm:"primaryKey"`
	Username    string
	Name        string
	Image       string
	CreatedDate time.Time
	DeletedDate sql.NullTime
}
