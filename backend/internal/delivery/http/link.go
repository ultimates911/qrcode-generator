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

	resp, err := h.linkUseCase.GetAllLinks(c.Context(), userID)
	if err != nil {
		c.Locals("logError", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

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