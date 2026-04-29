package handlers

import (
	"big-devops-api/internal/dto"
	"big-devops-api/internal/services"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	svc services.UserService
}

func NewUserHandler(svc services.UserService) *UserHandler {
	return &UserHandler{svc}
}

// GetUsers godoc
// @Summary List all users
// @Description Get a list of all users system-wide
// @Tags Users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} dto.UserResponse
// @Router /v1/users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.svc.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Register a new user with role and hashed password
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user body dto.UserRequest true "User Data"
// @Success 201 {object} dto.UserResponse
// @Router /v1/users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	req := new(dto.UserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := h.svc.Create(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUser godoc
// @Summary Get user by username
// @Description Get detailed information about a specific user
// @Tags Users
// @Security ApiKeyAuth
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} dto.UserResponse
// @Failure 404 {object} map[string]string
// @Router /v1/users/{username} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	username := c.Params("username")
	user, err := h.svc.GetByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return c.JSON(user)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Update user details such as name, email, or role
// @Tags Users
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Param user body dto.UserRequest true "Updated User Data"
// @Success 200 {object} dto.UserResponse
// @Router /v1/users/{username} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	username := c.Params("username")

	req := new(dto.UserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	user, err := h.svc.Update(username, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Remove a user from the system (Soft Delete)
// @Tags Users
// @Security ApiKeyAuth
// @Param username path string true "Username"
// @Success 204
// @Router /v1/users/{username} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	username := c.Params("username")
	if err := h.svc.Delete(username); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
