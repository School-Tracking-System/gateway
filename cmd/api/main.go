package main

import (
	"go.uber.org/fx"
)

// @title           Gateway Service API
// @version         1.0
// @description     API documentation for the gateway service of the School Tracking System.
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     JWT token in the format: Bearer {token}
func main() {
	fx.New(AppModule()).Run()
}
