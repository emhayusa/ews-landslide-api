package v1

import (
	"big-devops-api/internal/config"
	"big-devops-api/internal/handlers"
	"big-devops-api/internal/middleware"
	"big-devops-api/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(api fiber.Router, cfg *config.Config, authSvc services.AuthService, authHdl *handlers.AuthHandler, userHdl *handlers.UserHandler, stationHdl *handlers.StationHandler, streamHdl *handlers.StreamHandler) {
	v1 := api.Group("/v1")

	// WebSocket Streaming (Public or Protected depending on requirement, here public for testing)
	v1.Get("/ws", websocket.New(streamHdl.WSStream()))

	// Welcome & Health Check
	v1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to BIG-DEVOPS API V1",
			"status":  "Active",
			"version": "1.0.0",
		})
	})

	// Public Auth route (for local login)
	if cfg.AuthMethod == "local" {
		v1.Post("/login", authHdl.Login)
	}

	// Protected routes
	v1.Use(middleware.AuthMiddleware(authSvc))

	// Auth routes
	v1.Get("/me", authHdl.Me)
	v1.Get("/validate", authHdl.Validate)

	// User routes (Admin only)
	users := v1.Group("/users", middleware.RoleChecker("admin"))
	users.Get("/", userHdl.GetUsers)
	users.Post("/", userHdl.CreateUser)
	users.Get("/:username", userHdl.GetUser)
	users.Put("/:username", userHdl.UpdateUser)
	users.Delete("/:username", userHdl.DeleteUser)

	// Station routes
	stations := v1.Group("/stations")
	stations.Get("/", stationHdl.GetStations)
	stations.Post("/", stationHdl.CreateStation, middleware.RoleChecker("admin"))
	stations.Get("/:id", stationHdl.GetStation)
	stations.Put("/:id", stationHdl.UpdateStation, middleware.RoleChecker("admin"))
	stations.Delete("/:id", stationHdl.DeleteStation, middleware.RoleChecker("admin"))
}
