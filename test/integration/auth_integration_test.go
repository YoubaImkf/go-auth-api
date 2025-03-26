package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/YoubaImkf/go-auth-api/internal/controller"
	"github.com/YoubaImkf/go-auth-api/internal/dto"
	"github.com/YoubaImkf/go-auth-api/internal/middleware"
	"github.com/YoubaImkf/go-auth-api/internal/model"
	"github.com/YoubaImkf/go-auth-api/internal/repository"
	"github.com/YoubaImkf/go-auth-api/internal/service"
	"github.com/YoubaImkf/go-auth-api/test/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendPasswordResetEmail(to, token string) error {
	args := m.Called(to, token)
	return args.Error(0)
}

type AuthIntegrationTestSuite struct {
	suite.Suite
	db           *gorm.DB
	router       *gin.Engine
	config       *util.Config
	emailService *MockEmailService
}

func (suite *AuthIntegrationTestSuite) SetupSuite() {
	config, err := util.LoadTestConfig()
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.config = config

	db, err := gorm.Open("postgres", util.GetTestDSN(config))
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.db = db

	suite.emailService = new(MockEmailService)

	suite.db.AutoMigrate(&model.User{}, &model.PasswordReset{}, &model.BlacklistedToken{})

	suite.router = suite.setupTestRouter()
}

func (suite *AuthIntegrationTestSuite) setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())

	userRepo := repository.NewPostgresUserRepository(suite.db)
	blacklistRepo := repository.NewPostgresBlacklistRepository(suite.db)

	authService := service.NewAuthService(userRepo, blacklistRepo, suite.emailService)

	authController := controller.NewAuthController(authService)

	router.POST("/register", authController.Register)
	router.POST("/login", authController.Login)
	router.POST("/forgot-password", authController.ForgotPassword)
	router.POST("/reset-password", authController.ResetPassword)

	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(blacklistRepo))
	{
		protected.POST("/logout", authController.Logout)
		protected.GET("/me", authController.GetProfile)
	}

	return router
}

func (suite *AuthIntegrationTestSuite) SetupTest() {
	// Clean up database before each test
	suite.db.Exec("TRUNCATE users, password_resets, blacklisted_tokens RESTART IDENTITY CASCADE")

	// Reset mock expectations
	suite.emailService.ExpectedCalls = nil

	// Setup default mock behavior for email service
	suite.emailService.On("SendPasswordResetEmail", mock.Anything, mock.Anything).Return(nil)
}

func (suite *AuthIntegrationTestSuite) TestFullAuthFlow() {
	suite.emailService.On("SendPasswordResetEmail", "elon@example.com", mock.Anything).Return(nil)

	// 1. Register aa new user
	registerPayload := dto.RegisterRequest{
		Name:     "elon Musk",
		Email:    "elon@example.com",
		Password: "Password123!",
	}
	registerResp := suite.performRequest("POST", "/register", registerPayload)
	suite.Equal(http.StatusCreated, registerResp.Code)

	var registerResponse dto.RegisterResponse
	suite.NoError(json.Unmarshal(registerResp.Body.Bytes(), &registerResponse))
	suite.NotEmpty(registerResponse.AccessToken)

	// 2. Login with the registered user
	loginPayload := dto.LoginRequest{
		Email:    "elon@example.com",
		Password: "Password123!",
	}
	loginResp := suite.performRequest("POST", "/login", loginPayload)
	suite.Equal(http.StatusOK, loginResp.Code)

	var loginResponse dto.LoginResponse
	suite.NoError(json.Unmarshal(loginResp.Body.Bytes(), &loginResponse))
	suite.NotEmpty(loginResponse.AccessToken)

	// 3. Get user profile with the token
	profileResp := suite.performAuthorizedRequest("GET", "/me", nil, loginResponse.AccessToken)
	suite.Equal(http.StatusOK, profileResp.Code)

	var userResponse dto.UserResponse
	suite.NoError(json.Unmarshal(profileResp.Body.Bytes(), &userResponse))
	suite.Equal("elon Musk", userResponse.Name)

	// 4. Request password reset
	forgotPayload := dto.ForgotPasswordRequest{
		Email: "elon@example.com",
	}
	forgotResp := suite.performRequest("POST", "/forgot-password", forgotPayload)
	suite.Equal(http.StatusOK, forgotResp.Code)

	var forgotResponse map[string]interface{}
	suite.NoError(json.Unmarshal(forgotResp.Body.Bytes(), &forgotResponse))
	resetToken := forgotResponse["token"].(string)
	suite.NotEmpty(resetToken)

	// 5. Reset password
	resetPayload := dto.ResetPasswordRequest{
		Token:       resetToken,
		NewPassword: "NewPassword123!",
	}
	resetResp := suite.performRequest("POST", "/reset-password", resetPayload)
	suite.Equal(http.StatusNoContent, resetResp.Code)

	// 6. Login with new pasword
	loginPayload.Password = "NewPassword123!"
	loginResp = suite.performRequest("POST", "/login", loginPayload)
	suite.Equal(http.StatusOK, loginResp.Code)

	// 7. Logout
	logoutResp := suite.performAuthorizedRequest("POST", "/logout", nil, loginResponse.AccessToken)
	suite.Equal(http.StatusNoContent, logoutResp.Code)

	// 8. Verify token is blacklisted 🤞
	profileResp = suite.performAuthorizedRequest("GET", "/me", nil, loginResponse.AccessToken)
	suite.Equal(http.StatusUnauthorized, profileResp.Code)

	suite.emailService.AssertExpectations(suite.T())
}

func (suite *AuthIntegrationTestSuite) TestInvalidLoginAttempt() {
	registerPayload := dto.RegisterRequest{
		Name:     "elon Musk",
		Email:    "elon@example.com",
		Password: "Password123!",
	}
	suite.performRequest("POST", "/register", registerPayload)

	loginPayload := dto.LoginRequest{
		Email:    "elon@example.com",
		Password: "WrongPassword123!",
	}
	loginResp := suite.performRequest("POST", "/login", loginPayload)

	suite.Equal(http.StatusUnauthorized, loginResp.Code)
}

func (suite *AuthIntegrationTestSuite) TestProtectedRouteWithoutAuth() {
	resp := suite.performRequest("GET", "/me", nil)
	suite.Equal(http.StatusUnauthorized, resp.Code)
}

func (suite *AuthIntegrationTestSuite) TestInvalidToken() {
	resp := suite.performAuthorizedRequest("GET", "/me", nil, "invalid_token")
	suite.Equal(http.StatusUnauthorized, resp.Code)
}

func (suite *AuthIntegrationTestSuite) TestResetPasswordTokenReuse() {
	registerPayload := dto.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "Password123!",
	}
	suite.performRequest("POST", "/register", registerPayload)

	forgotPayload := dto.ForgotPasswordRequest{
		Email: "test@example.com",
	}
	forgotResp := suite.performRequest("POST", "/forgot-password", forgotPayload)

	var forgotResponse map[string]any
	// parse the JSON response from the forgot password API into the 'forgotResponse' map.
	// 'forgotResp.Body.Bytes()' gets the response body as a byte slice.
	suite.NoError(json.Unmarshal(forgotResp.Body.Bytes(), &forgotResponse))
	resetToken := forgotResponse["token"].(string)

	resetPayload := dto.ResetPasswordRequest{
		Token:       resetToken,
		NewPassword: "NewPassword123!",
	}
	firstResetResp := suite.performRequest("POST", "/reset-password", resetPayload)
	suite.Equal(http.StatusNoContent, firstResetResp.Code)

	// Second attempt with same token should fail 🥹 !!!
	secondResetResp := suite.performRequest("POST", "/reset-password", resetPayload)
	suite.Equal(http.StatusUnauthorized, secondResetResp.Code)
}

// --- Pirvate Method ---
func (suite *AuthIntegrationTestSuite) performRequest(method, path string, payload interface{}) *httptest.ResponseRecorder {
	var body []byte
	if payload != nil {
		body, _ = json.Marshal(payload)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	return w
}

func (suite *AuthIntegrationTestSuite) performAuthorizedRequest(method, path string, payload interface{}, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	if payload != nil {
		body, _ := json.Marshal(payload)
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	return w
}

// --- End Pirvate Method ---

func (suite *AuthIntegrationTestSuite) TearDownSuite() {
	suite.db.Close()
}

func TestAuthIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(AuthIntegrationTestSuite))
}
