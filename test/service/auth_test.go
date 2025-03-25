package service

import (
	"fmt"
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
	err := godotenv.Load("../../.env")
	if err != nil {
		suite.T().Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DATABASE_HOST")
	if dbHost == "db" {
		dbHost = "localhost" // Use localhost for local testing
	}

	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	testDbName := "go-auth-db-test"

	db, err := gorm.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		dbHost, dbPort, dbUser, testDbName, dbPassword))
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.db = db

	suite.db.AutoMigrate(&model.User{}, &model.PasswordReset{})

	suite.userRepo = repository.NewPostgresUserRepository(suite.db)
	suite.blacklistRepo = repository.NewPostgresBlacklistRepository(suite.db)
	suite.emailService = new(MockEmailService)
	suite.authService = service.NewAuthService(suite.userRepo, suite.blacklistRepo, suite.emailService)
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
	assert.NotEmpty(suite.T(), accessToken)
	assert.NotEmpty(suite.T(), refreshToken)
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

	user, accessToken, refreshToken, err := suite.authService.Login(loginRequest)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.NotEmpty(suite.T(), accessToken)
	assert.NotEmpty(suite.T(), refreshToken)
}

func (suite *AuthServiceTestSuite) TestForgotPassword() {
	user := &model.User{
		Name:  "johndoe",
		Email: "john.doe@example.com",
	}
	suite.userRepo.Create(user)

	suite.emailService.On("SendPasswordResetEmail", "john.doe@example.com", mock.Anything).Return(nil)

	err := suite.authService.ForgotPassword("john.doe@example.com")

	assert.NoError(suite.T(), err)
	suite.emailService.AssertExpectations(suite.T())
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
