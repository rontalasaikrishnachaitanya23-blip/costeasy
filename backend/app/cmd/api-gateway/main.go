package main

import (
	"log"

	"github.com/chaitu35/costeasy/backend/app/cmd/api-gateway/bootstrap"
	_ "github.com/chaitu35/costeasy/backend/docs" // ✅ Swagger generated files
	_ "github.com/chaitu35/costeasy/backend/internal/payroll/attendance/handler"
	
)

// @title CostEasy API Gateway
// @version 1.0
// @description API Gateway for CostEasy backend modules (Auth, GL, Payroll, Settings)
// @termsOfService http://costeasy.io/terms/

// @contact.name API Support
// @contact.email support@costeasy.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http
func main() {
	app, err := bootstrap.InitializeApp()
	if err != nil {
		log.Fatalf("bootstrap failed: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
