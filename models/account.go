package models

import (
	"database/sql"
	"time"
)

// Account is a creator account that has a profile on the platform
type Account struct {
	ID          uint64 `gorm:"primaryKey"`
	Email       string
	Secret      string
	CreatedDate time.Time
	DeletedDate sql.NullTime
}
