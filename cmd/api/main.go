package main

import (
	_ "github.com/YoubaImkf/go-auth-api/docs"
	"github.com/YoubaImkf/go-auth-api/internal/app"
)

// @title           Authentication API
// @version         1.0
// @description     A JWT Authentication Service API
// @host            localhost:8080
// @BasePath        /03622bf7-d58b-4997-965c-14ee58c63554/
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	application := app.New()
	application.Run()
}
