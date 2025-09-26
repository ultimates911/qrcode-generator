package http

import (
	"qrcodegen/internal/dto"
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
