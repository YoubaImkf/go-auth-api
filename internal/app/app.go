package app

import (
	"fmt"
	"log"
	"os"

	"github.com/YoubaImkf/go-auth-api/internal/controller"
	"github.com/YoubaImkf/go-auth-api/internal/middleware"
	"github.com/YoubaImkf/go-auth-api/internal/model"
	"github.com/YoubaImkf/go-auth-api/internal/repository"
	"github.com/YoubaImkf/go-auth-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	router *gin.Engine
	db     *gorm.DB
}

func New() *App {
	app := &App{
		router: gin.Default(),
	}

	app.loadEnv()
	app.loadConfig()
	app.initDB()
	app.setupRoutes()
	return app
}

func (a *App) loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func (a *App) loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func (a *App) initDB() {
	config := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"),
	)

	db, err := gorm.Open("postgres", config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}
	a.db = db

	// Auto Migrate the User model an PasswordReset to create the tables
	if err := a.db.AutoMigrate(&model.User{}, &model.PasswordReset{}).Error; err != nil {
		log.Fatalf("Failed to auto-migrate models: %s", err)
	}
}

func (a *App) setupRoutes() {
	blacklistRepo := repository.NewPostgresBlacklistRepository(a.db)
	userRepo := repository.NewPostgresUserRepository(a.db)
	authService := service.NewAuthService(userRepo, blacklistRepo)

	healthController := controller.NewHealthController()
	authController := controller.NewAuthController(authService)

	a.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	a.router.GET("/api/v1/health", healthController.Health)

	// Auth routes
	a.router.POST("/api/v1/auth/register", authController.Register)
	a.router.POST("/api/v1/auth/login", authController.Login)
	a.router.POST("/api/v1/auth/forgot-password", authController.ForgotPassword)
	a.router.POST("/api/v1/auth/reset-password", authController.ResetPassword)

	// Needed : Apply AuthMidleware only to the logout route !
	protected := a.router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(blacklistRepo))
	protected.POST("/auth/logout", authController.Logout)
	protected.GET("/auth/me", authController.GetProfile)
}
func (a *App) Run() {
	a.router.Run(":8080")
}
