package main

// @title BIG-DEVOPS API
// @version 1.0
// @description API for DEVOPS.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

import (
	_ "big-devops-api/docs" // Import generated docs
	"big-devops-api/internal/config"
	"big-devops-api/internal/database"
	"big-devops-api/internal/handlers"
	"big-devops-api/internal/repositories"
	v1 "big-devops-api/internal/routes/v1"
	"big-devops-api/internal/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Initialize Database
	database.ConnectDB(cfg)
	database.ConnectRedis(cfg)

	// Repositories
	userRepo := repositories.NewUserRepository(database.DB)
	stationRepo := repositories.NewStationRepository(database.DB)
	monitoringRepo := repositories.NewMonitoringRepository(database.DB)

	// Services
	authSvc := services.NewAuthService(cfg, userRepo)
	userSvc := services.NewUserService(userRepo)
	stationSvc := services.NewStationService(stationRepo)
	streamSvc := services.NewStreamService(cfg, monitoringRepo, stationRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authSvc, userRepo)
	userHandler := handlers.NewUserHandler(userSvc)
	stationHandler := handlers.NewStationHandler(stationSvc)
	streamHandler := handlers.NewStreamHandler(streamSvc)

	app := fiber.New(fiber.Config{
		AppName: "BIG-DEVOPS API",
	})

	// Global Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Setup Routes
	api := app.Group("/api")
	v1.SetupRoutes(api, cfg, authSvc, authHandler, userHandler, stationHandler, streamHandler)

	// Swagger UI
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Start server
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
