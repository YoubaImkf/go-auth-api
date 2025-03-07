package model

import "time"

type PasswordReset struct {
	Email  string    `gorm:"primary_key"`
	Token  string    `gorm:"unique"`
	Expiry time.Time `gorm:"index"`
}
