package service

import (
	"testing"
	"time"

	"github.com/YoubaImkf/go-auth-api/internal/dto"
	"github.com/YoubaImkf/go-auth-api/internal/model"
	"github.com/YoubaImkf/go-auth-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

// RRemoveAll implements repository.UserRepository.
func (m *MockUserRepository) RemoveAll() error {
	panic("unimplemented")
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByUserNameOrEmail(identifier string) (*model.User, error) {
	args := m.Called(identifier)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) StorePasswordResetToken(email, token string, expiry time.Time) error {
	args := m.Called(email, token, expiry)
	return args.Error(0)
}

func (m *MockUserRepository) FindEmailByResetToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) UpdatePassword(email, newPassword string) error {
	args := m.Called(email, newPassword)
	return args.Error(0)
}

type MockBlacklistRepository struct {
	mock.Mock
}

func (m *MockBlacklistRepository) Add(token string, expiry time.Time) error {
	args := m.Called(token, expiry)
	return args.Error(0)
}

func (m *MockBlacklistRepository) IsBlacklisted(token string) bool {
	args := m.Called(token)
	return args.Bool(0)
}

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendPasswordResetEmail(to, token string) error {
	args := m.Called(to, token)
	return args.Error(0)
}

// User service
func (m *MockUserRepository) GetAll() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func TestAuthService_Register(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBlacklistRepo := new(MockBlacklistRepository)
	mockEmailService := new(MockEmailService)
	authService := service.NewAuthService(mockUserRepo, mockBlacklistRepo, mockEmailService)

	registerRequest := dto.RegisterRequest{
		Name:     "johndoe",
		Email:    "john.doe@example.com",
		Password: "password123",
	}

	mockUserRepo.On("Create", mock.Anything).Return(nil)

	user, accessToken, refreshToken, err := authService.Register(registerRequest)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBlacklistRepo := new(MockBlacklistRepository)
	mockEmailService := new(MockEmailService)
	authService := service.NewAuthService(mockUserRepo, mockBlacklistRepo, mockEmailService)

	loginRequest := dto.LoginRequest{
		Email:    "johndoe",
		Password: "password123",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUser := &model.User{
		Name:     "johndoe",
		Email:    "john.doe@example.com",
		Password: string(hashedPassword),
	}

	mockUserRepo.On("FindByEmail", "john.doe@example.com").Return(mockUser, nil)

	user, accessToken, refreshToken, err := authService.Login(loginRequest)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
	mockUserRepo.AssertExpectations(t)
}

// func TestAuthService_Logout(t *testing.T) {
//     mockUserRepo := new(MockUserRepository)
//     mockBlacklistRepo := new(MockBlacklistRepository)
//     mockEmailService := new(MockEmailService)
//     authService := service.NewAuthService(mockUserRepo, mockBlacklistRepo, mockEmailService)

//     token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//         "sub": "john.doe@example.com",
//         "exp": time.Now().Add(time.Hour).Unix(),
//     })
//     tokenString, _ := token.SignedString([]byte("your_jwt_secret"))

//     mockBlacklistRepo.On("Add", tokenString, mock.Anything).Return(nil)

//     err := authService.Logout(tokenString)

//     assert.NoError(t, err)
//     mockBlacklistRepo.AssertExpectations(t)
// }

func TestAuthService_ForgotPassword(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBlacklistRepo := new(MockBlacklistRepository)
	mockEmailService := new(MockEmailService)
	authService := service.NewAuthService(mockUserRepo, mockBlacklistRepo, mockEmailService)

	email := "john.doe@example.com"
	mockUser := &model.User{
		Name:  "johndoe",
		Email: email,
	}

	mockUserRepo.On("FindByEmail", email).Return(mockUser, nil)
	mockUserRepo.On("StorePasswordResetToken", email, mock.Anything, mock.Anything).Return(nil)
	mockEmailService.On("SendPasswordResetEmail", email, mock.Anything).Return(nil)

	err := authService.ForgotPassword(email)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

func TestAuthService_ResetPassword(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockBlacklistRepo := new(MockBlacklistRepository)
	mockEmailService := new(MockEmailService)
	authService := service.NewAuthService(mockUserRepo, mockBlacklistRepo, mockEmailService)

	resetPasswordRequest := dto.ResetPasswordRequest{
		Token:       "reset_token",
		NewPassword: "newpassword123",
	}

	email := "john.doe@example.com"
	mockUserRepo.On("FindEmailByResetToken", "reset_token").Return(email, nil)
	mockUserRepo.On("UpdatePassword", email, mock.Anything).Return(nil)

	err := authService.ResetPassword(resetPasswordRequest)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

// User service
func TestUserService_GetAllUsers(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	userService := service.NewUserService(mockUserRepo)

	mockUsers := []model.User{
		{
			Name:  "johndoe",
			Email: "john.doe@example.com",
		},
		{
			Name:  "johndoe",
			Email: "jane.doe@example.com",
		},
	}

	mockUserRepo.On("GetAll").Return(mockUsers, nil)

	users, err := userService.GetAllUsers()

	assert.NoError(t, err)
	assert.Equal(t, mockUsers, users)
	mockUserRepo.AssertExpectations(t)
}
