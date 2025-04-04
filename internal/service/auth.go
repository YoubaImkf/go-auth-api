package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/YoubaImkf/go-auth-api/internal/dto"
	"github.com/YoubaImkf/go-auth-api/internal/model"
	"github.com/YoubaImkf/go-auth-api/internal/repository"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// const (
// 	minPasswordLength = 8
// )

type AuthService struct {
	userRepository repository.UserRepository
	blacklistRepo  repository.BlacklistRepository
	emailService   EmailService
	jwtSecret      string
}

func NewAuthService(userRepo repository.UserRepository, blacklistRepo repository.BlacklistRepository, emailService EmailService) *AuthService {
	return &AuthService{
		userRepository: userRepo,
		blacklistRepo:  blacklistRepo,
		emailService:   emailService,
		jwtSecret:      viper.GetString("jwt.secret"),
	}
}

func (s *AuthService) Register(registerRequest dto.RegisterRequest) (*model.User, string, string, error) {
	// if err := isPasswordValid(registerRequest.Password); err != nil {
	// 	return nil, "", "", err
	// }

	existingUser, err := s.userRepository.FindByEmail(registerRequest.Email)
	if err == nil && existingUser != nil {
		return nil, "", "", errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", err
	}

	user := &model.User{
		Name:     registerRequest.Name,
		Email:    registerRequest.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepository.Create(user); err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *AuthService) Login(loginRequest dto.LoginRequest) (*model.User, string, string, error) {
	user, err := s.userRepository.FindByEmail(loginRequest.Email)
	if err != nil {
		return nil, "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	accessToken, refreshToken, err := s.generateTokens(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *AuthService) Logout(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiry := time.Unix(int64(claims["exp"].(float64)), 0)
		return s.blacklistRepo.Add(tokenString, expiry)
	}

	return errors.New("invalid token")
}

func (s *AuthService) GetUserProfile(email string) (*model.User, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) ForgotPassword(email string) (string, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	resetToken := hex.EncodeToString(token)

	// Store the token :3
	expiry := time.Now().Add(1 * time.Hour)
	if err := s.userRepository.StorePasswordResetToken(user.Email, resetToken, expiry); err != nil {
		return "", err
	}
	// Send the token via email
	if err := s.emailService.SendPasswordResetEmail(user.Email, resetToken); err != nil {
		return "", err
	}

	log.Printf("Password reset token for %s: %s", user.Email, resetToken)

	return resetToken, nil
}

func (s *AuthService) ResetPassword(resetPasswordRequest dto.ResetPasswordRequest) error {
	email, err := s.userRepository.FindEmailByResetToken(resetPasswordRequest.Token)
	if err != nil {
		return errors.New("invalid or expired reset token")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(resetPasswordRequest.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.userRepository.UpdatePassword(email, string(hashedPassword)); err != nil {
		return err
	}

	if err := s.userRepository.InvalidateResetToken(resetPasswordRequest.Token); err != nil {
		return err
	}

	return nil
}

// --- Private Methods ---

func (s *AuthService) generateTokens(user *model.User) (string, string, error) {
	accessTokenExpiry := time.Now().Add(viper.GetDuration("jwt.access_token_expiry"))
	refreshTokenExpiry := time.Now().Add(viper.GetDuration("jwt.refresh_token_expiry"))

	accessTokenClaims := jwt.MapClaims{
		"sub": user.Email,
		"exp": accessTokenExpiry.Unix(),
	}

	refreshTokenClaims := jwt.MapClaims{
		"sub": user.Email,
		"exp": refreshTokenExpiry.Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// func isPasswordValid(password string) error {
// 	if len(password) < minPasswordLength {
// 		return errors.New("password must be at least 8 characters long")
// 	}

// 	var (
// 		hasUpper   bool
// 		hasLower   bool
// 		hasNumber  bool
// 		hasSpecial bool
// 	)

// 	for _, char := range password {
// 		switch {
// 		case unicode.IsUpper(char):
// 			hasUpper = true
// 		case unicode.IsLower(char):
// 			hasLower = true
// 		case unicode.IsNumber(char):
// 			hasNumber = true
// 		case unicode.IsPunct(char) || unicode.IsSymbol(char):
// 			hasSpecial = true
// 		}
// 	}

// 	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
// 		return errors.New("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
// 	}

// 	return nil
// }
