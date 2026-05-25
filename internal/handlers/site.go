package handlers

import (
	"big-devops-api/internal/dto"
	"big-devops-api/internal/services"

	"github.com/gofiber/fiber/v2"
)

type SiteHandler struct {
	svc services.SiteService
}

func NewSiteHandler(svc services.SiteService) *SiteHandler {
	return &SiteHandler{svc}
}

// GetSites godoc
// @Summary List all sites
// @Description Get a list of all observation sites
// @Tags Sites
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} dto.SiteResponse
// @Router /v1/sites [get]
func (h *SiteHandler) GetSites(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	role := c.Locals("role").(string)

	sites, err := h.svc.GetAll(userID, role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(sites)
}

// CreateSite godoc
// @Summary Create a new site
// @Description Register a new observation site
// @Tags Sites
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param site body dto.SiteRequest true "Site Data"
// @Success 201 {object} dto.SiteResponse
// @Router /v1/sites [post]
func (h *SiteHandler) CreateSite(c *fiber.Ctx) error {
	req := new(dto.SiteRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	site, err := h.svc.Create(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(site)
}

// GetSite godoc
// @Summary Get site by ID
// @Description Get detailed information about a specific site
// @Tags Sites
// @Security ApiKeyAuth
// @Produce json
// @Param id path string true "Site ID"
// @Success 200 {object} dto.SiteResponse
// @Failure 404 {object} map[string]string
// @Router /v1/sites/{id} [get]
func (h *SiteHandler) GetSite(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("user_id").(uint)
	role := c.Locals("role").(string)

	site, err := h.svc.GetByID(id, userID, role)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Site not found or access denied"})
	}
	return c.JSON(site)
}

// UpdateSite godoc
// @Summary Update an existing site
// @Description Update site details
// @Tags Sites
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Site ID"
// @Param site body dto.SiteRequest true "Updated Site Data"
// @Success 200 {object} dto.SiteResponse
// @Router /v1/sites/{id} [put]
func (h *SiteHandler) UpdateSite(c *fiber.Ctx) error {
	id := c.Params("id")

	req := new(dto.SiteRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	site, err := h.svc.Update(id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(site)
}

// DeleteSite godoc
// @Summary Delete a site
// @Description Remove a site from the system (Soft Delete)
// @Tags Sites
// @Security ApiKeyAuth
// @Param id path string true "Site ID"
// @Success 204
// @Router /v1/sites/{id} [delete]
func (h *SiteHandler) DeleteSite(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
