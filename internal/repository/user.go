package repository

import (
	"errors"
	"time"

	"github.com/YoubaImkf/go-auth-api/internal/model"
	"github.com/jinzhu/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByUserNameOrEmail(identifier string) (*model.User, error)
	StorePasswordResetToken(email, token string, expiry time.Time) error
	FindEmailByResetToken(token string) (string, error)
	UpdatePassword(email, newPassword string) error
	GetAll() ([]model.User, error)
	RemoveAll() error
}

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) Create(user *model.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostgresUserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *PostgresUserRepository) FindByUserNameOrEmail(identifier string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("user_name = ? OR email = ?", identifier, identifier).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *PostgresUserRepository) StorePasswordResetToken(email, token string, expiry time.Time) error {
	passwordReset := model.PasswordReset{
		Email:  email,
		Token:  token,
		Expiry: expiry,
	}
	return r.db.Save(&passwordReset).Error
}

func (r *PostgresUserRepository) FindEmailByResetToken(token string) (string, error) {
	var passwordReset model.PasswordReset
	if err := r.db.Where("token = ? AND expiry > ?", token, time.Now()).First(&passwordReset).Error; err != nil {
		return "", errors.New("invalid or expired reset token")
	}
	return passwordReset.Email, nil
}

func (r *PostgresUserRepository) UpdatePassword(email, newPassword string) error {
	return r.db.Model(&model.User{}).Where("email = ?", email).Update("password", newPassword).Error
}

func (r *PostgresUserRepository) GetAll() ([]model.User, error) {
	var users []model.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *PostgresUserRepository) RemoveAll() error {
	if err := r.db.Delete(&model.User{}).Error; err != nil {
		return err
	}
	return nil
}
