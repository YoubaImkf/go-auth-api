package main

import (
	_ "github.com/YoubaImkf/go-auth-api/docs"
	"github.com/YoubaImkf/go-auth-api/internal/app"
)

// @title           Authentication API
// @version         1.0
// @description     A JWT Authentication Service API
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	application := app.New()
	application.Run()
}
