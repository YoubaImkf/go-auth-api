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
	if os.Getenv("APP_ENVIRONMENT") != "production" {
		// In development, load environment variables from .env
		if err := godotenv.Load(); err != nil {
			log.Printf("Error loading .env file: %v", err)
		}
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

	// Bind specific environment variables to override config values
	if err := viper.BindEnv("database.host", "DATABASE_HOST"); err != nil {
		log.Printf("Error binding DATABASE_HOST: %v", err)
	}
	if err := viper.BindEnv("database.port", "DATABASE_PORT"); err != nil {
		log.Printf("Error binding SMTP_HOST: %v", err)
	}

	if err := viper.BindEnv("smtp.host", "SMTP_HOST"); err != nil {
		log.Printf("Error binding SMTP_HOST: %v", err)
	}
	if err := viper.BindEnv("smtp.port", "SMTP_PORT"); err != nil {
		log.Printf("Error binding DATABASE_HOST: %v", err)
	}
	if err := viper.BindEnv("smtp.username", "SMTP_USERNAME"); err != nil {
		log.Printf("Error binding DATABASE_HOST: %v", err)
	}
	if err := viper.BindEnv("smtp.password", "SMTP_PASSWORD"); err != nil {
		log.Printf("Error binding DATABASE_HOST: %v", err)
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
	groupUUID := viper.GetString("group.uuid")

	blacklistRepo := repository.NewPostgresBlacklistRepository(a.db)
	userRepo := repository.NewPostgresUserRepository(a.db)
	emailService := service.NewEmailService()
	authService := service.NewAuthService(userRepo, blacklistRepo, emailService)
	userService := service.NewUserService(userRepo)

	healthController := controller.NewHealthController()
	authController := controller.NewAuthController(authService)
	userController := controller.NewUserController(userService)

	a.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Group routes under the UUID
	apiGroup := a.router.Group(groupUUID)

	apiGroup.GET("/health", healthController.Health)
	apiGroup.GET("/users", userController.GetAllUsers)

	authGroup := apiGroup.Group("/auth")
	authGroup.POST("/register", authController.Register)
	authGroup.POST("/login", authController.Login)
	authGroup.POST("/forgot-password", authController.ForgotPassword)
	authGroup.POST("/reset-password", authController.ResetPassword)

	protected := apiGroup.Group("/")
	protected.Use(middleware.AuthMiddleware(blacklistRepo))
	protected.POST("/auth/logout", authController.Logout)
	protected.GET("/auth/me", authController.GetProfile)
}

func (a *App) Run() {
	a.router.Run(":8080")
}
