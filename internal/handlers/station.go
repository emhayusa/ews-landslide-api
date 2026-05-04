package handlers

import (
	"big-devops-api/internal/dto"
	"big-devops-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type StationHandler struct {
	svc services.StationService
}

func NewStationHandler(svc services.StationService) *StationHandler {
	return &StationHandler{svc}
}

// GetStations godoc
// @Summary List all stations
// @Description Get a list of all EWS Landslide stations
// @Tags Stations
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} dto.StationResponse
// @Router /v1/stations [get]
func (h *StationHandler) GetStations(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	role := c.Locals("role").(string)

	stations, err := h.svc.GetAll(userID, role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(stations)
}

// CreateStation godoc
// @Summary Create a new station
// @Description Register a new EWS Landslide station
// @Tags Stations
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param station body dto.StationRequest true "Station Data"
// @Success 201 {object} dto.StationResponse
// @Router /v1/stations [post]
func (h *StationHandler) CreateStation(c *fiber.Ctx) error {
	req := new(dto.StationRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	station, err := h.svc.Create(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(station)
}

// GetStation godoc
// @Summary Get station by ID
// @Description Get detailed information about a specific station
// @Tags Stations
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Station ID"
// @Success 200 {object} dto.StationResponse
// @Failure 404 {object} map[string]string
// @Router /v1/stations/{id} [get]
func (h *StationHandler) GetStation(c *fiber.Ctx) error {
	id := c.Params("id")
	station, err := h.svc.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Station not found"})
	}
	return c.JSON(station)
}

// UpdateStation godoc
// @Summary Update an existing station
// @Description Update station details
// @Tags Stations
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Station ID"
// @Param station body dto.StationRequest true "Updated Station Data"
// @Success 200 {object} dto.StationResponse
// @Router /v1/stations/{id} [put]
func (h *StationHandler) UpdateStation(c *fiber.Ctx) error {
	id := c.Params("id")

	req := new(dto.StationRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	station, err := h.svc.Update(id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(station)
}

// DeleteStation godoc
// @Summary Delete a station
// @Description Remove a station from the system (Soft Delete)
// @Tags Stations
// @Security ApiKeyAuth
// @Param id path string true "Station ID"
// @Success 204
// @Router /v1/stations/{id} [delete]
func (h *StationHandler) DeleteStation(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
