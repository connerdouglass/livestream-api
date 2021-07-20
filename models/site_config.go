package models

import (
	"database/sql"
)

// SiteConfig is some general configuration data for the site overall
type SiteConfig struct {
	ID              uint64 `gorm:"primaryKey"`
	VapidPublicKey  sql.NullString
	VapidPrivateKey sql.NullString
}
