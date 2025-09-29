package http

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"qrcodegen/internal/dto"
	"qrcodegen/internal/pkg/qrcode"
	"qrcodegen/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type QRHandler struct {
	validate  *validator.Validate
	qrUseCase *usecase.QRUseCase
}

func NewQRHandler(validate *validator.Validate, qrUseCase *usecase.QRUseCase) *QRHandler {
	return &QRHandler{validate: validate, qrUseCase: qrUseCase}
}

func (h *QRHandler) Generate(c *fiber.Ctx) error {
	var req dto.GenerateQRCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}
	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	png, err := h.qrUseCase.Generate(c.Context(), req)
	if err != nil {
		c.Locals("logError", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate qr-code"})
	}

	c.Type("png")
	return c.Status(fiber.StatusCreated).Send(png)
}

func (h *LinkHandler) DownloadQR(c *fiber.Ctx) error {
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

	link, err := h.linkUseCase.GetLinkByID(c.Context(), int64(linkID), userID)
	if err != nil {
		if errors.Is(err, usecase.ErrLinkNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		c.Locals("logError", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	typ := strings.ToLower(c.Query("type", "png"))
	if typ == "" {
		typ = "png"
	}

	var smoothing float64
	if link.Smoothing != nil {
		smoothing = *link.Smoothing
	}

	var (
		data        []byte
		contentType string
		filename    string
	)
	switch typ {
	case "png":
		data, err = qrcode.GeneratePNG(link.OriginalURL, link.Color, link.Background, smoothing)
		contentType = "image/png"
		filename = fmt.Sprintf("qr-%d.png", linkID)
	case "svg":
		data, err = qrcode.GenerateSVG(link.OriginalURL, link.Color, link.Background, smoothing)
		contentType = "image/svg+xml"
		filename = fmt.Sprintf("qr-%d.svg", linkID)
	case "pdf":
		data, err = qrcode.GeneratePDF(link.OriginalURL, link.Color, link.Background, smoothing)
		contentType = "application/pdf"
		filename = fmt.Sprintf("qr-%d.pdf", linkID)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "query param 'type' must be one of: png, svg, pdf"})
	}
	if err != nil {
		c.Locals("logError", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate qr-code"})
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	return c.Status(fiber.StatusOK).Send(data)
}