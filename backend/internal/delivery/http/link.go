package http

import (
	"errors"
	"strconv"

	"qrcodegen/config"

	"qrcodegen/internal/dto"
	"qrcodegen/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type LinkHandler struct {
	validate    *validator.Validate
	linkUseCase *usecase.LinkUseCase
	cfg         *config.Config
}

func NewLinkHandler(validate *validator.Validate, linkUseCase *usecase.LinkUseCase, cfg *config.Config) *LinkHandler {
	return &LinkHandler{
		validate:    validate,
		linkUseCase: linkUseCase,
		cfg:         cfg,
	}
}

// CreateLink godoc
// @Summary Create a new link
// @Description Create a new shortened link for the authenticated user
// @Tags links
// @Accept  json
// @Produce  json
// @Param   link  body      dto.CreateLinkRequest  true  "Link data"
// @Success 201   {object}  dto.CreateLinkResponse
// @Failure 400   {object}  dto.GenericError
// @Failure 401   {object}  dto.GenericError
// @Failure 500   {object}  dto.GenericError
// @Router /links/create [post]
func (h *LinkHandler) CreateLink(c *fiber.Ctx) error {
	var req dto.CreateLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.Locals("logError", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	resp, err := h.linkUseCase.CreateLink(c.Context(), req, userID)
	if err != nil {
		c.Locals("logError", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// GetAllLinks godoc
// @Summary Get all links for a user
// @Description Get all links created by the authenticated user
// @Tags links
// @Produce  json
// @Param   search  query  string  false  "Filter by link name (case-insensitive)"
// @Param   sort_by  query  string  false  "Sort by: created_at|transitions"
// @Param   order    query  string  false  "Sort order: asc|desc"
// @Success 200 {object} dto.GetAllLinksResponse
// @Failure 401 {object} dto.GenericError
// @Failure 500 {object} dto.GenericError
// @Router /links [get]
func (h *LinkHandler) GetAllLinks(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.Locals("logError", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	search := c.Query("search")
	sortBy := c.Query("sort_by")
	order := c.Query("order")
	var resp *dto.GetAllLinksResponse
	if search != "" {
		resp, err = h.linkUseCase.SearchLinksByName(c.Context(), userID, search)
	} else {
		resp, err = h.linkUseCase.GetAllLinks(c.Context(), userID)
	}
	if err != nil {
		c.Locals("logError", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	// apply sorting in usecase layer helper
	by := usecase.SortByCreatedAt
	if sortBy == "transitions" {
		by = usecase.SortByTransitions
	}
	ord := usecase.SortDesc
	if order == "asc" {
		ord = usecase.SortAsc
	}
	resp.Links = h.linkUseCase.SortLinks(resp.Links, by, ord)

	return c.Status(fiber.StatusOK).JSON(resp)
}

// GetLink godoc
// @Summary Get a link by ID
// @Description Get a specific link by its ID for the authenticated user
// @Tags links
// @Produce  json
// @Param   id   path      int  true  "Link ID"
// @Success 200 {object} dto.GetLinkResponse
// @Failure 400 {object} dto.GenericError
// @Failure 401 {object} dto.GenericError
// @Failure 404 {object} dto.GenericError
// @Failure 500 {object} dto.GenericError
// @Router /links/{id} [get]
func (h *LinkHandler) GetLink(c *fiber.Ctx) error {
	linkID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid link ID"})
	}

	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.Locals("logError", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	resp, err := h.linkUseCase.GetLinkByID(c.Context(), int64(linkID), userID)
	if err != nil {
		if errors.Is(err, usecase.ErrLinkNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		c.Locals("logError", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// EditLink godoc
// @Summary Edit a link
// @Description Edit a specific link by its ID for the authenticated user
// @Tags links
// @Accept  json
// @Produce  json
// @Param   id   path      int  true  "Link ID"
// @Param   link body      dto.EditLinkRequest  true  "Updated link data"
// @Success 200 {object} dto.EditLinkResponse
// @Failure 400 {object} dto.GenericError
// @Failure 401 {object} dto.GenericError
// @Failure 404 {object} dto.GenericError
// @Failure 500 {object} dto.GenericError
// @Router /links/{id} [patch]
func (h *LinkHandler) EditLink(c *fiber.Ctx) error {
	linkID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid link ID"})
	}

	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.Locals("logError", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	var req dto.EditLinkRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp, err := h.linkUseCase.EditLink(c.Context(), int64(linkID), userID, req)
	if err != nil {
		if errors.Is(err, usecase.ErrLinkNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		c.Locals("logError", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// Redirect godoc
// @Summary Redirect to original URL
// @Description Redirects a shortened link to its original URL
// @Tags redirect
// @Param   hash   path      string  true  "Link hash"
// @Success 302 {string} string "Redirects to the original URL"
// @Failure 400 {object} dto.GenericError
// @Failure 404 {object} dto.GenericError
// @Failure 500 {object} dto.GenericError
// @Router /redirect/{hash} [get]
func (h *LinkHandler) Redirect(c *fiber.Ctx) error {
	hash := c.Params("hash")
	if hash == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Hash is required"})
	}

	referer := c.Get("Referer")
	userAgent := c.Get("User-Agent")
	ip := c.IP()

	originalURL, err := h.linkUseCase.Redirect(c.Context(), hash, referer, userAgent, ip)
	if err != nil {
		if errors.Is(err, usecase.ErrLinkNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Redirect(originalURL, fiber.StatusFound)
}

// GetTransitionsByLink godoc
// @Summary Get transitions for a link
// @Description Get transition analytics for a specific link by its ID
// @Tags links
// @Produce  json
// @Param   id   path      int  true  "Link ID"
// @Success 200 {object} dto.GetTransitionsResponse
// @Failure 400 {object} dto.GenericError
// @Failure 401 {object} dto.GenericError
// @Failure 500 {object} dto.GenericError
// @Router /links/{id}/transitions [get]
func (h *LinkHandler) GetTransitionsByLink(c *fiber.Ctx) error {
	linkID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid link ID"})
	}
	userIDStr, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.Locals("logError", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	resp, err := h.linkUseCase.GetTransitions(c.Context(), int64(linkID), userID)
	if err != nil {
		c.Locals("logError", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
