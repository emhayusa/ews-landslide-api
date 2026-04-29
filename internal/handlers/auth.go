package handlers

import (
	"big-devops-api/internal/dto"
	"big-devops-api/internal/repositories"
	"big-devops-api/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type AuthHandler struct {
	authService services.AuthService
	userRepo    repositories.UserRepository
}

func NewAuthHandler(authService services.AuthService, userRepo repositories.UserRepository) *AuthHandler {
	return &AuthHandler{authService, userRepo}
}

// Login godoc
// @Summary Authenticate user
// @Description Authenticate user via username and password (local) or proxy to Keycloak
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body dto.LoginRequest true "Login Credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /v1/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(dto.LoginResponse{Token: token})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	claims := c.Locals("user").(jwt.MapClaims)
	
	sub := ""
	if s, ok := claims["sub"].(string); ok {
		sub = s
	}

	user, err := h.userRepo.FindByUsername(sub)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(dto.UserMeResponse{
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	})
}

func (h *AuthHandler) Validate(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "valid",
		"valid":  true,
	})
}
