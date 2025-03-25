package app

import (
	"fmt"
	"log"
	"os"

	"github.com/YoubaImkf/go-auth-api/docs"
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
	app.InitSwaggerHost()
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
}

func (a *App) initDB() {
	config := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		viper.GetString("DATABASE_HOST"),
		viper.GetInt("DATABASE_PORT"),
		viper.GetString("DATABASE_USER"),
		viper.GetString("DATABASE_NAME"),
		viper.GetString("DATABASE_PASSWORD"),
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

func (a *App) InitSwaggerHost() {
	if host := os.Getenv("APP_HOST"); host != "" {
		docs.SwaggerInfo.Host = host
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

	apiGroup.POST("/register", authController.Register)
	apiGroup.POST("/login", authController.Login)
	apiGroup.POST("/forgot-password", authController.ForgotPassword)
	apiGroup.POST("/reset-password", authController.ResetPassword)

	protected := apiGroup.Group("/")
	protected.Use(middleware.AuthMiddleware(blacklistRepo))
	protected.POST("/logout", authController.Logout)
	protected.GET("/me", authController.GetProfile)
}

func (a *App) Run() {
	a.router.Run(":8080")
}
