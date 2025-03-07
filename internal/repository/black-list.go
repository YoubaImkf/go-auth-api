package repository

import (
	"time"

	"github.com/jinzhu/gorm"
)

type BlacklistRepository interface {
	Add(token string, expiry time.Time) error
	IsBlacklisted(token string) bool
}

type PostgresBlacklistRepository struct {
	db *gorm.DB
}

type BlacklistedToken struct {
	Token  string    `gorm:"primary_key"`
	Expiry time.Time `gorm:"index"`
}

func NewPostgresBlacklistRepository(db *gorm.DB) *PostgresBlacklistRepository {
	db.AutoMigrate(&BlacklistedToken{})
	return &PostgresBlacklistRepository{db: db}
}

func (r *PostgresBlacklistRepository) Add(token string, expiry time.Time) error {
	blacklistedToken := BlacklistedToken{
		Token:  token,
		Expiry: expiry,
	}
	return r.db.Create(&blacklistedToken).Error
}

func (r *PostgresBlacklistRepository) IsBlacklisted(token string) bool {
	var blacklistedToken BlacklistedToken
	if err := r.db.Where("token = ?", token).First(&blacklistedToken).Error; err != nil {
		return false
	}
	return time.Now().Before(blacklistedToken.Expiry)
}
