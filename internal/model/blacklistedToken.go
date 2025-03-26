package model

import "time"

type BlacklistedToken struct {
	Token  string    `gorm:"primary_key;type:text"`
	Expiry time.Time `gorm:"index;not null"`
}
