package controller

import (
	"net/http"
	"strings"

	"github.com/YoubaImkf/go-auth-api/internal/dto"
	"github.com/YoubaImkf/go-auth-api/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// @Summary      Register user
// @Description  Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  dto.RegisterRequest  true  "User"
// @Success      201  {object}  dto.RegisterResponse
// @Router       /register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var registerRequest dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&registerRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, accessToken, refreshToken, err := c.authService.Register(registerRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.RegisterResponse{
		User: dto.UserResponse{
			Name:  user.Name,
			Email: user.Email,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	ctx.JSON(http.StatusCreated, response)
}

// @Summary      Login user
// @Description  Authenticate a user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  dto.LoginRequest  true  "User"
// @Success      200  {object}  dto.LoginResponse
// @Router       /login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var loginRequest dto.LoginRequest

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, accessToken, refreshToken, err := c.authService.Login(loginRequest)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	response := dto.LoginResponse{
		User: dto.UserResponse{
			Name:  user.Name,
			Email: user.Email,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary      Logout user
// @Description  Logout a user
// @Tags         auth
// @Produce      json
// @Success      204  {object}  map[string]interface{}
// @Router       /logout [post]
// @Security     Bearer
func (c *AuthController) Logout(ctx *gin.Context) {
	tokenString := ctx.GetHeader("Authorization")
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if tokenString == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
		return
	}

	err := c.authService.Logout(tokenString)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "Successfully logged out"})
}

// @Summary      Get user profile
// @Description  Get the profile of the logged-in user
// @Tags         auth
// @Produce      json
// @Success      200  {object}  dto.UserResponse
// @Router       /me [get]
// @Security     Bearer
func (c *AuthController) GetProfile(ctx *gin.Context) {
	userEmail, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	user, err := c.authService.GetUserProfile(userEmail.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.UserResponse{
		Name:  user.Name,
		Email: user.Email,
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary      Forgot password
// @Description  Request a password reset
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        email  body  dto.ForgotPasswordRequest  true  "Email"
// @Success      204  {object}  map[string]interface{}
// @Router       /forgot-password [post]
func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var forgotPasswordRequest dto.ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&forgotPasswordRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.authService.ForgotPassword(forgotPasswordRequest.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "Password reset link sent", "token": token})
}

// @Summary      Reset password
// @Description  Reset the user's password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        resetPasswordRequest  body  dto.ResetPasswordRequest  true  "Reset Password"
// @Success      200  {object}  map[string]interface{}
// @Router       /reset-password [post]
func (c *AuthController) ResetPassword(ctx *gin.Context) {
	var resetPasswordRequest dto.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&resetPasswordRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.authService.ResetPassword(resetPasswordRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password has been reset"})
}
