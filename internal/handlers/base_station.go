package handlers

import (
	"big-devops-api/internal/dto"
	"big-devops-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type BaseStationHandler struct {
	svc services.BaseStationService
}

func NewBaseStationHandler(svc services.BaseStationService) *BaseStationHandler {
	return &BaseStationHandler{svc}
}

// GetBaseStations godoc
// @Summary List all base stations
// @Description Get a list of all Base Stations
// @Tags BaseStations
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} dto.BaseStationResponse
// @Router /v1/base-stations [get]
func (h *BaseStationHandler) GetBaseStations(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	role := c.Locals("role").(string)

	stations, err := h.svc.GetAll(userID, role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stations)
}

// CreateBaseStation godoc
// @Summary Create a new base station
// @Description Register a new Base Station
// @Tags BaseStations
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param station body dto.BaseStationRequest true "Base Station Data"
// @Success 201 {object} dto.BaseStationResponse
// @Router /v1/base-stations [post]
func (h *BaseStationHandler) CreateBaseStation(c *fiber.Ctx) error {
	req := new(dto.BaseStationRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	station, err := h.svc.Create(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(station)
}

// GetBaseStation godoc
// @Summary Get base station by ID or UUID
// @Description Get detailed information about a specific base station
// @Tags BaseStations
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Base Station ID or UUID"
// @Success 200 {object} dto.BaseStationResponse
// @Failure 404 {object} map[string]string
// @Router /v1/base-stations/{id} [get]
func (h *BaseStationHandler) GetBaseStation(c *fiber.Ctx) error {
	id := c.Params("id")
	station, err := h.svc.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Base Station not found"})
	}
	return c.JSON(station)
}

// UpdateBaseStation godoc
// @Summary Update an existing base station
// @Description Update base station details
// @Tags BaseStations
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Base Station ID"
// @Param station body dto.BaseStationRequest true "Updated Base Station Data"
// @Success 200 {object} dto.BaseStationResponse
// @Router /v1/base-stations/{id} [put]
func (h *BaseStationHandler) UpdateBaseStation(c *fiber.Ctx) error {
	id := c.Params("id")

	req := new(dto.BaseStationRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	station, err := h.svc.Update(id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(station)
}

// DeleteBaseStation godoc
// @Summary Delete a base station
// @Description Remove a base station from the system (Soft Delete)
// @Tags BaseStations
// @Security ApiKeyAuth
// @Param id path string true "Base Station ID"
// @Success 204
// @Router /v1/base-stations/{id} [delete]
func (h *BaseStationHandler) DeleteBaseStation(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
