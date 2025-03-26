package service

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/YoubaImkf/go-auth-api/internal/dto"
	"github.com/YoubaImkf/go-auth-api/internal/model"
	"github.com/YoubaImkf/go-auth-api/internal/repository"
	"github.com/YoubaImkf/go-auth-api/internal/service"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendPasswordResetEmail(to, token string) error {
	args := m.Called(to, token)
	return args.Error(0)
}

type AuthServiceTestSuite struct {
	suite.Suite
	db            *gorm.DB
	userRepo      repository.UserRepository
	blacklistRepo repository.BlacklistRepository
	emailService  *MockEmailService
	authService   *service.AuthService
}

func (suite *AuthServiceTestSuite) SetupSuite() {
	suite.loadEnvironmentVariables()
	suite.initializeDatabase()
	suite.initializeRepositories()
	suite.initializeServices()
}

func (suite *AuthServiceTestSuite) TearDownSuite() {
	suite.db.Close()
}

func (suite *AuthServiceTestSuite) SetupTest() {
	// Clean up the database before each test
	suite.db.Exec("DELETE FROM users")
	suite.db.Exec("DELETE FROM password_resets")
}

func (suite *AuthServiceTestSuite) TestRegister() {
	registerRequest := dto.RegisterRequest{
		Name:     "johndoe",
		Email:    "john.doe@example.com",
		Password: "Password123!",
	}

	user, accessToken, refreshToken, err := suite.authService.Register(registerRequest)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), registerRequest.Name, user.Name)
	assert.Equal(suite.T(), registerRequest.Email, user.Email)
	assert.NotEmpty(suite.T(), user.Password)
	assert.NotEqual(suite.T(), registerRequest.Password, user.Password)
	assert.NotEmpty(suite.T(), accessToken)
	assert.NotEmpty(suite.T(), refreshToken)

	var savedUser model.User
	err = suite.db.Where("email = ?", registerRequest.Email).First(&savedUser).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Name, savedUser.Name)
	assert.Equal(suite.T(), user.Email, savedUser.Email)
}

func (suite *AuthServiceTestSuite) TestRegisterExistingUser() {
	existingUser := &model.User{
		Name:     "johndoe",
		Email:    "john.doe@example.com",
		Password: "hashedpassword",
	}
	suite.userRepo.Create(existingUser)

	registerRequest := dto.RegisterRequest{
		Name:     "johndoe",
		Email:    "john.doe@example.com",
		Password: "Password123!",
	}

	user, accessToken, refreshToken, err := suite.authService.Register(registerRequest)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user already exists", err.Error())
	assert.Nil(suite.T(), user)
	assert.Empty(suite.T(), accessToken)
	assert.Empty(suite.T(), refreshToken)
}

func (suite *AuthServiceTestSuite) TestRegisterInvalidPassword() {
	testCases := []struct {
		name     string
		password string
		wantErr  string
	}{
		{
			name:     "too short password",
			password: "weak",
			wantErr:  "password must be at least 8 characters long",
		},
		// {
		// 	name:     "missing uppercase",
		// 	password: "password123!",
		// 	wantErr:  "password must contain at least one uppercase letter, one lowercase letter, one number, and one special character",
		// },
		// {
		// 	name:     "missing lowercase",
		// 	password: "PASSWORD123!",
		// 	wantErr:  "password must contain at least one uppercase letter, one lowercase letter, one number, and one special character",
		// },
		// {
		// 	name:     "missing number",
		// 	password: "Password!",
		// 	wantErr:  "password must contain at least one uppercase letter, one lowercase letter, one number, and one special character",
		// },
		// {
		// 	name:     "missing special character",
		// 	password: "Password123",
		// 	wantErr:  "password must contain at least one uppercase letter, one lowercase letter, one number, and one special character",
		// },
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			registerRequest := dto.RegisterRequest{
				Name:     "johndoe",
				Email:    "john.doe@example.com",
				Password: tc.password,
			}

			user, accessToken, refreshToken, err := suite.authService.Register(registerRequest)

			assert.Error(t, err)
			assert.Equal(t, tc.wantErr, err.Error())
			assert.Nil(t, user)
			assert.Empty(t, accessToken)
			assert.Empty(t, refreshToken)
		})
	}
}

func (suite *AuthServiceTestSuite) TestLogin() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
	user := &model.User{
		Name:     "johndoe",
		Email:    "john.doe@example.com",
		Password: string(hashedPassword),
	}
	suite.userRepo.Create(user)

	loginRequest := dto.LoginRequest{
		Email:    "john.doe@example.com",
		Password: "Password123!",
	}

	loggedInUser, accessToken, refreshToken, err := suite.authService.Login(loginRequest)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), loggedInUser)
	assert.Equal(suite.T(), user.Email, loggedInUser.Email)
	assert.Equal(suite.T(), user.Name, loggedInUser.Name)
	assert.NotEmpty(suite.T(), accessToken)
	assert.NotEmpty(suite.T(), refreshToken)
}

func (suite *AuthServiceTestSuite) TestLoginInvalidCredentials() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)
	user := &model.User{
		Name:     "johndoe",
		Email:    "john.doe@example.com",
		Password: string(hashedPassword),
	}
	suite.userRepo.Create(user)

	loginRequest := dto.LoginRequest{
		Email:    "john.doe@example.com",
		Password: "WrongPassword123!",
	}

	loggedInUser, accessToken, refreshToken, err := suite.authService.Login(loginRequest)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "invalid credentials", err.Error())
	assert.Nil(suite.T(), loggedInUser)
	assert.Empty(suite.T(), accessToken)
	assert.Empty(suite.T(), refreshToken)
}

func (suite *AuthServiceTestSuite) TestLoginNonExistentUser() {
	loginRequest := dto.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "Password123!",
	}

	loggedInUser, accessToken, refreshToken, err := suite.authService.Login(loginRequest)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), loggedInUser)
	assert.Empty(suite.T(), accessToken)
	assert.Empty(suite.T(), refreshToken)
}
func (suite *AuthServiceTestSuite) TestForgotPassword() {
	// Setup
	user := &model.User{
		Name:  "johndoe",
		Email: "john.doe@example.com",
	}
	suite.userRepo.Create(user)

	suite.emailService.On("SendPasswordResetEmail", "john.doe@example.com", mock.Anything).Return(nil)

	token, err := suite.authService.ForgotPassword("john.doe@example.com")

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
	assert.Len(suite.T(), token, 64)
	suite.emailService.AssertExpectations(suite.T())

	var passwordReset model.PasswordReset
	err = suite.db.Where("email = ? AND token = ?", user.Email, token).First(&passwordReset).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Email, passwordReset.Email)
	assert.Equal(suite.T(), token, passwordReset.Token)
}

func (suite *AuthServiceTestSuite) TestForgotPasswordUserNotFound() {
	token, err := suite.authService.ForgotPassword("nonexistent@example.com")

	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), token)
	assert.Equal(suite.T(), "user not found", err.Error())
	suite.emailService.AssertNotCalled(suite.T(), "SendPasswordResetEmail")
}

func (suite *AuthServiceTestSuite) TestResetPassword() {
	user := &model.User{
		Name:  "johndoe",
		Email: "john.doe@example.com",
	}
	suite.userRepo.Create(user)

	token := "reset_token"
	expiry := time.Now().Add(1 * time.Hour)
	suite.userRepo.StorePasswordResetToken(user.Email, token, expiry)

	resetPasswordRequest := dto.ResetPasswordRequest{
		Token:       token,
		NewPassword: "newPassword123!",
	}

	err := suite.authService.ResetPassword(resetPasswordRequest)

	assert.NoError(suite.T(), err)
}

func (suite *AuthServiceTestSuite) TestGetAllUsers() {
	user1 := &model.User{
		Name:  "johndoe",
		Email: "john.doe@example.com",
	}
	user2 := &model.User{
		Name:  "janedoe",
		Email: "jane.doe@example.com",
	}
	suite.userRepo.Create(user1)
	suite.userRepo.Create(user2)

	users, err := suite.userRepo.GetAll()

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 2)
}

func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

// --- Private Methods ---

func (suite *AuthServiceTestSuite) loadEnvironmentVariables() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
}

func (suite *AuthServiceTestSuite) getDBConfig() (host, port, user, password string) {
	host = suite.getEnvOrDefault("DATABASE_HOST", "localhost")
	port = suite.getEnvOrDefault("DATABASE_PORT", "5432")
	user = suite.getEnvOrDefault("POSTGRES_USER", "root")
	password = suite.getEnvOrDefault("POSTGRES_PASSWORD", "lets-jungle-it-bro!")
	return
}

func (suite *AuthServiceTestSuite) getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (suite *AuthServiceTestSuite) initializeDatabase() {
	host, port, user, password := suite.getDBConfig()
	testDbName := "go-auth-db-test"

	db, err := suite.connectToDatabase(host, port, user, password, testDbName)
	if err != nil {
		suite.T().Fatalf("Failed to connect to database: %v", err)
	}

	suite.db = db
	suite.migrateDatabase()
}

func (suite *AuthServiceTestSuite) connectToDatabase(host, port, user, password, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbName, password,
	)
	return gorm.Open("postgres", dsn)
}

func (suite *AuthServiceTestSuite) migrateDatabase() {
	suite.db.AutoMigrate(&model.User{}, &model.PasswordReset{})
}

func (suite *AuthServiceTestSuite) initializeRepositories() {
	suite.userRepo = repository.NewPostgresUserRepository(suite.db)
	suite.blacklistRepo = repository.NewPostgresBlacklistRepository(suite.db)
}

func (suite *AuthServiceTestSuite) initializeServices() {
	suite.emailService = new(MockEmailService)
	suite.authService = service.NewAuthService(
		suite.userRepo,
		suite.blacklistRepo,
		suite.emailService,
	)
}
